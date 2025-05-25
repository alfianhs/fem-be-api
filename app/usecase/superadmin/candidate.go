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

	// extract seasonTeamIds
	seasonTeamIds := helpers.ExtractIds(candidates, func(c mongo_model.Candidate) string {
		return c.SeasonTeamID
	})

	// fetch seasonTeams
	seasonTeams, err := u.mongoDbRepo.FetchListSeasonTeam(ctx, map[string]interface{}{"ids": seasonTeamIds})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer seasonTeams.Close(ctx)

	seasonTeamsMap := make(map[string]mongo_model.SeasonTeam)
	for seasonTeams.Next(ctx) {
		var s mongo_model.SeasonTeam
		if err := seasonTeams.Decode(&s); err != nil {
			logrus.Error("SeasonTeam Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		seasonTeamsMap[s.ID.Hex()] = s
	}

	// extract seasonTeamPlayerIds
	seasonTeamPlayerIds := helpers.ExtractIds(candidates, func(c mongo_model.Candidate) string {
		return c.SeasonTeamPlayerID
	})

	// fetch seasonTeamPlayers
	seasonTeamPlayers, err := u.mongoDbRepo.FetchListSeasonTeamPlayer(ctx, map[string]interface{}{"ids": seasonTeamPlayerIds})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer seasonTeamPlayers.Close(ctx)

	seasonTeamPlayersMap := make(map[string]mongo_model.SeasonTeamPlayer)
	for seasonTeamPlayers.Next(ctx) {
		var s mongo_model.SeasonTeamPlayer
		if err := seasonTeamPlayers.Decode(&s); err != nil {
			logrus.Error("SeasonTeamPlayer Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		seasonTeamPlayersMap[s.ID.Hex()] = s
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

		st, ok := seasonTeamsMap[c.SeasonTeamID]
		if ok {
			c.SeasonTeam = mongo_model.SeasonTeamFK{
				ID:       st.ID.Hex(),
				SeasonID: st.SeasonID,
				TeamID:   st.Team.ID,
				Team:     st.Team,
			}
		} else {
			c.SeasonTeam = mongo_model.SeasonTeamFK{}
			logrus.Error("SeasonTeam not found:", c.SeasonTeamID)
		}

		stp, ok := seasonTeamPlayersMap[c.SeasonTeamPlayerID]
		if ok {
			c.SeasonTeamPlayer = mongo_model.SeasonTeamPlayerFK{
				ID:         stp.ID.Hex(),
				SeasonTeam: stp.SeasonTeam,
				Player:     stp.Player,
				Position:   stp.Position,
				Image:      stp.Image.URL,
			}
		} else {
			c.SeasonTeamPlayer = mongo_model.SeasonTeamPlayerFK{}
			logrus.Error("SeasonTeamPlayer not found:", c.SeasonTeamPlayerID)
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
	}
	c.Voting = mongo_model.VotingFK{
		ID:         voting.ID.Hex(),
		Title:      voting.Title,
		TotalVoter: voting.TotalVoter,
	}

	// fetch seasonTeam
	seasonTeam, err := u.mongoDbRepo.FetchOneSeasonTeam(ctx, map[string]interface{}{"id": c.SeasonTeamID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeam == nil {
		c.SeasonTeam = mongo_model.SeasonTeamFK{}
		logrus.Error("SeasonTeam not found:", c.SeasonTeamID)
	}
	c.SeasonTeam = mongo_model.SeasonTeamFK{
		ID:       seasonTeam.ID.Hex(),
		SeasonID: seasonTeam.SeasonID,
		TeamID:   seasonTeam.Team.ID,
		Team:     seasonTeam.Team,
	}

	// fetch seasonTeamPlayer
	seasonTeamPlayer, err := u.mongoDbRepo.FetchOneSeasonTeamPlayer(ctx, map[string]interface{}{"id": c.SeasonTeamPlayerID})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if seasonTeamPlayer == nil {
		c.SeasonTeamPlayer = mongo_model.SeasonTeamPlayerFK{}
		logrus.Error("SeasonTeamPlayer not found:", c.SeasonTeamPlayerID)
	}
	c.SeasonTeamPlayer = mongo_model.SeasonTeamPlayerFK{
		ID:         seasonTeamPlayer.ID.Hex(),
		SeasonTeam: seasonTeamPlayer.SeasonTeam,
		Player:     seasonTeamPlayer.Player,
		Position:   seasonTeamPlayer.Position,
		Image:      seasonTeamPlayer.Image.URL,
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

	now := time.Now()
	candidate := &mongo_model.Candidate{
		ID:                 primitive.NewObjectID(),
		VotingID:           voting.ID.Hex(),
		SeasonTeamID:       seasonTeamPlayer.SeasonTeam.ID,
		SeasonTeamPlayerID: seasonTeamPlayer.ID.Hex(),
		Performance:        payload.Performance,
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
		c.SeasonTeamPlayerID = seasonTeamPlayer.ID.Hex()
		fields["seasonTeamPlayerId"] = c.SeasonTeamPlayerID
	}
	if payload.Performance != "" {
		c.Performance = payload.Performance
		fields["performance"] = c.Performance
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
