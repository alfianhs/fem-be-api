package superadmin_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCandidateList
func (u *superadminAppUsecase) GetCandidateList(ctx context.Context, queryParam url.Values) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, offset, limit := helpers.GetOffsetLimit(queryParam)
	opts := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	if v := queryParam.Get("votingId"); v != "" {
		opts["votingId"] = v
	}
	if stp := queryParam.Get("seasonTeamPlayerId"); stp != "" {
		opts["seasonTeamPlayerId"] = stp
	}

	total := u.mongoDbRepo.CountCandidate(ctx, opts)
	if total == 0 {
		return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
			List:  []interface{}{},
			Limit: limit,
			Page:  page,
			Total: total,
		})
	}

	if sort := queryParam.Get("sort"); sort != "" {
		opts["sort"] = sort
	}
	if dir := queryParam.Get("dir"); dir != "" {
		opts["dir"] = dir
	}

	cur, err := u.mongoDbRepo.FetchListCandidate(ctx, opts)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer cur.Close(ctx)

	var candidates []mongo_model.Candidate
	for cur.Next(ctx) {
		var c mongo_model.Candidate
		if err := cur.Decode(&c); err != nil {
			logrus.Error("Candidate Decode Error:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		candidates = append(candidates, c)
	}

	// extract votingIds
	votingIds := helpers.ExtractIds(candidates, func(c mongo_model.Candidate) string {
		return c.VotingID
	})

	// fetch votings
	votings, err := u.mongoDbRepo.FetchListVoting(ctx, map[string]interface{}{"ids": votingIds})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer votings.Close(ctx)

	votingsMap := make(map[string]mongo_model.Voting)
	for votings.Next(ctx) {
		var s mongo_model.Voting
		if err := votings.Decode(&s); err != nil {
			logrus.Error("Voting Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		votingsMap[s.ID.Hex()] = s
	}

	var list []interface{}
	for _, c := range candidates {
		v, ok := votingsMap[c.VotingID]
		if ok {
			c.Voting = mongo_model.VotingFK{
				ID:         v.ID.Hex(),
				Title:      v.Title,
				TotalVoter: v.TotalVoter,
			}
		} else {
			c.Voting = mongo_model.VotingFK{}
			logrus.Error("Voting not found:", c.VotingID)
		}

		list = append(list, c.Format(&c.Voting))
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		List:  list,
		Limit: limit,
		Page:  page,
		Total: total,
	})
}

// GetCandidateDetail
func (u *superadminAppUsecase) GetCandidateDetail(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	c, err := u.mongoDbRepo.FetchOneCandidate(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if c == nil {
		return helpers.NewResponse(http.StatusNotFound, "Candidate not found", nil, nil)
	}

	// fetch voting
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{"id": c.VotingID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting == nil {
		c.Voting = mongo_model.VotingFK{}
		logrus.Error("Voting not found:", c.VotingID)
	} else {
		c.Voting = mongo_model.VotingFK{
			ID:         voting.ID.Hex(),
			Title:      voting.Title,
			TotalVoter: voting.TotalVoter,
		}
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, c.Format(&c.Voting))
}

// CreateCandidate
func (u *superadminAppUsecase) CreateCandidate(ctx context.Context, payload request.CandidateCreateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1) Validate payload
	errs := map[string]string{}
	if payload.VotingID == "" {
		errs["votingId"] = "Voting ID is required"
	}
	if payload.SeasonTeamPlayerID == "" {
		errs["seasonTeamPlayerId"] = "SeasonTeamPlayer ID is required"
	}
	if payload.Performance.TeamLeaderboard == 0 {
		errs["performance.teamLeaderboard"] = "TeamLeaderboard is required"
	}
	if payload.Performance.Goal == 0 && payload.Performance.Assist == 0 && payload.Performance.Save == 0 {
		errs["performance.goal"] = "At least one of goal, assist, or save must be provided"
		errs["performance.assist"] = "At least one of goal, assist, or save must be provided"
		errs["performance.save"] = "At least one of goal, assist, or save must be provided"
	}

	if len(errs) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation error", errs, nil)
	}

	// 2) Ensure voting exists
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{"id": payload.VotingID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Voting not found", nil, nil)
	}

	// ensure seasonTeamPlayer exists
	seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{"id": payload.SeasonTeamPlayerID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeamPlayer == nil {
		return helpers.NewResponse(http.StatusBadRequest, "SeasonTeamPlayer not found", nil, nil)
	}

	// get seasonTeam
	seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{"id": seasonTeamPlayer.SeasonTeam.ID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeam == nil {
		return helpers.NewResponse(http.StatusBadRequest, "SeasonTeam not found", nil, nil)
	}

	// score point from voting
	scorePoint := helpers.PerformancePoint{
		Goal:   voting.PerformancePoint.Goal,
		Assist: voting.PerformancePoint.Assist,
		Save:   voting.PerformancePoint.Save,
	}

	// performance count from payload
	candidatePerformanceCount := helpers.CandidatePerformanceCount{
		Goal:   payload.Performance.Goal,
		Assist: payload.Performance.Assist,
		Save:   payload.Performance.Save,
	}

	now := time.Now()
	candidate := &mongo_model.Candidate{
		ID:       primitive.NewObjectID(),
		VotingID: voting.ID.Hex(),
		SeasonTeam: mongo_model.SeasonTeamFK{
			ID:       seasonTeam.ID.Hex(),
			SeasonID: seasonTeam.SeasonID,
			TeamID:   seasonTeam.Team.ID,
			Team:     seasonTeam.Team,
		},
		SeasonTeamPlayer: mongo_model.SeasonTeamPlayerFK{
			ID:         seasonTeamPlayer.ID.Hex(),
			SeasonTeam: seasonTeamPlayer.SeasonTeam,
			Player:     seasonTeamPlayer.Player,
			Position:   seasonTeamPlayer.Position,
			Image:      seasonTeamPlayer.Image.URL,
		},
		Performance: mongo_model.Performance{
			TeamLeaderboard: payload.Performance.TeamLeaderboard,
			Goal:            payload.Performance.Goal,
			Assist:          payload.Performance.Assist,
			Save:            payload.Performance.Save,
			Score:           helpers.CalculateScore(scorePoint, candidatePerformanceCount),
		},
		Voters: mongo_model.Voters{
			Count: 0,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.mongoDbRepo.CreateOneCandidate(ctx, candidate); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusCreated, "Candidate created", nil, candidate)
}

// UpdateCandidate
func (u *superadminAppUsecase) UpdateCandidate(ctx context.Context, id string, payload request.CandidateUpdateRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1) Fetch existing candidate

	c, err := u.mongoDbRepo.FetchOneCandidate(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if c == nil {
		return helpers.NewResponse(http.StatusNotFound, "Candidate not found", nil, nil)
	}

	// 2) Prevent edit if voting is active
	now := time.Now()
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{"id": c.VotingID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting != nil && voting.Status == mongo_model.VotingStatusActive &&
		now.After(voting.StartDate) && now.Before(voting.EndDate) {
		return helpers.NewResponse(http.StatusBadRequest, "Cannot edit candidate while its voting is active", nil, nil)
	}

	// 3) Apply updates
	fields := map[string]interface{}{}
	if payload.SeasonTeamPlayerID != "" {
		// ensure seasonTeamPlayer exists
		seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{"id": payload.SeasonTeamPlayerID})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		if seasonTeamPlayer == nil {
			return helpers.NewResponse(http.StatusBadRequest, "SeasonTeamPlayer not found", nil, nil)
		}

		// check season team
		seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{"id": seasonTeamPlayer.SeasonTeam.ID})
		if err != nil {
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		if seasonTeam == nil {
			return helpers.NewResponse(http.StatusBadRequest, "SeasonTeam not found", nil, nil)
		}
		c.SeasonTeam = mongo_model.SeasonTeamFK{
			ID:       seasonTeam.ID.Hex(),
			SeasonID: seasonTeam.SeasonID,
			TeamID:   seasonTeam.Team.ID,
			Team:     seasonTeam.Team,
		}
		fields["seasonTeam"] = c.SeasonTeam
		c.SeasonTeamPlayer = mongo_model.SeasonTeamPlayerFK{
			ID:         seasonTeamPlayer.ID.Hex(),
			SeasonTeam: seasonTeamPlayer.SeasonTeam,
			Player:     seasonTeamPlayer.Player,
			Position:   seasonTeamPlayer.Position,
			Image:      seasonTeamPlayer.Image.URL,
		}
		fields["seasonTeamPlayer"] = c.SeasonTeamPlayer
	}

	c.Performance.TeamLeaderboard = payload.Performance.TeamLeaderboard
	fields["performance.teamLeaderboard"] = c.Performance.TeamLeaderboard

	c.Performance.Goal = payload.Performance.Goal
	fields["performance.goal"] = c.Performance.Goal

	c.Performance.Assist = payload.Performance.Assist
	fields["performance.assist"] = c.Performance.Assist

	c.Performance.Save = payload.Performance.Save
	fields["performance.save"] = c.Performance.Save

	if payload.Performance.Goal != 0 || payload.Performance.Assist != 0 || payload.Performance.Save != 0 {
		// score point from voting
		scorePoint := helpers.PerformancePoint{
			Goal:   voting.PerformancePoint.Goal,
			Assist: voting.PerformancePoint.Assist,
			Save:   voting.PerformancePoint.Save,
		}

		// performance count from updated candidate
		candidatePerformanceCount := helpers.CandidatePerformanceCount{
			Goal:   c.Performance.Goal,
			Assist: c.Performance.Assist,
			Save:   c.Performance.Save,
		}
		fields["performance.score"] = helpers.CalculateScore(scorePoint, candidatePerformanceCount)
	} else {
		errs := map[string]string{}
		errs["performance.goal"] = "At least one of goal, assist, or save must be provided"
		errs["performance.assist"] = "At least one of goal, assist, or save must be provided"
		errs["performance.save"] = "At least one of goal, assist, or save must be provided"

		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation error", errs, nil)
	}

	if len(fields) == 0 {
		return helpers.NewResponse(http.StatusBadRequest, "No fields to update", nil, nil)
	} else {
		c.UpdatedAt = now
		fields["updatedAt"] = c.UpdatedAt
	}

	if err := u.mongoDbRepo.UpdatePartialCandidate(ctx, map[string]interface{}{"id": c.ID}, fields); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Candidate updated", nil, c)
}

// DeleteCandidate
func (u *superadminAppUsecase) DeleteCandidate(ctx context.Context, id string) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1) Fetch existing candidate
	c, err := u.mongoDbRepo.FetchOneCandidate(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if c == nil {
		return helpers.NewResponse(http.StatusNotFound, "Candidate not found", nil, nil)
	}

	// 2) Prevent delete if voting is active
	now := time.Now()
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{"id": c.VotingID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting != nil && voting.Status == mongo_model.VotingStatusActive &&
		now.After(voting.StartDate) && now.Before(voting.EndDate) {
		return helpers.NewResponse(http.StatusBadRequest, "Cannot delete candidate while its voting is active", nil, nil)
	}

	// 3) Soft-delete
	if err := u.mongoDbRepo.UpdatePartialCandidate(ctx, map[string]interface{}{"id": id}, map[string]interface{}{
		"deletedAt": now,
	}); err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Candidate deleted", nil, nil)
}
