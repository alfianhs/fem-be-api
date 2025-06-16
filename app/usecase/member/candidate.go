package member_usecase

import (
	mongo_model "app/domain/model/mongo"
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *memberAppUsecase) GetCandidateList(ctx context.Context, claim jwt_helpers.MemberJWTClaims, queryParam url.Values) helpers.Response {
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

	// extract candidateIds
	candidateIds := helpers.ExtractIds(candidates, func(c mongo_model.Candidate) string {
		return c.ID.Hex()
	})

	// fetch voting logs
	optionsVotingLogs := map[string]interface{}{
		"candidate.ids": candidateIds,
		"memberId":      claim.UserID,
	}
	votingLogs, err := u.mongoDbRepo.FetchListVotingLog(ctx, optionsVotingLogs)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer votingLogs.Close(ctx)

	// map chosen candidate
	chosenCandidateMap := make(map[string]bool)
	for votingLogs.Next(ctx) {
		var s mongo_model.VotingLog
		if err := votingLogs.Decode(&s); err != nil {
			logrus.Error("VotingLog Decode:", err)
			return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
		}
		chosenCandidateMap[s.Candidate.ID] = true
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

		if chosenCandidateMap[c.ID.Hex()] {
			isChosen := true
			c.IsChosen = &isChosen
		}

		list = append(list, c.Format(&c.Voting))
	}

	// sorting
	sort.SliceStable(list, func(i, j int) bool {
		a := list[i].(*mongo_model.Candidate)
		b := list[j].(*mongo_model.Candidate)

		// compare by voters percentage
		if a.Voters.Percentage != b.Voters.Percentage {
			return a.Voters.Percentage > b.Voters.Percentage
		}

		// compare by candidate performance score
		return a.Performance.Score > b.Performance.Score
	})

	return helpers.NewResponse(http.StatusOK, "Success", nil, helpers.PaginatedResponse{
		List:  list,
		Limit: limit,
		Page:  page,
		Total: total,
	})
}

func (u *memberAppUsecase) CandidateVote(ctx context.Context, claim jwt_helpers.MemberJWTClaims, payload request.CandidateVoteRequest) helpers.Response {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check member
	member, err := u.mongoDbRepo.FetchOneMember(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if member == nil {
		return helpers.NewResponse(http.StatusBadRequest, "User not found", nil, nil)
	}

	// validate payload
	errValidation := make(map[string]string)
	if payload.CandidateID == "" {
		errValidation["candidateId"] = "Candidate ID field is required"
	}
	if payload.VotingID == "" {
		errValidation["votingId"] = "Voting ID field is required"
	}
	if len(errValidation) > 0 {
		return helpers.NewResponse(http.StatusUnprocessableEntity, "Validation Error", errValidation, nil)
	}

	// check voting
	voting, err := u.mongoDbRepo.FetchOneVoting(ctx, map[string]interface{}{
		"id": payload.VotingID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if voting == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Voting not found", nil, nil)
	}

	// check candidate
	candidate, err := u.mongoDbRepo.FetchOneCandidate(ctx, map[string]interface{}{
		"id": payload.CandidateID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if candidate == nil {
		return helpers.NewResponse(http.StatusBadRequest, "Candidate not found", nil, nil)
	}

	// validate voting status
	now := time.Now()
	if voting.Status != mongo_model.VotingStatusActive {
		return helpers.NewResponse(http.StatusBadRequest, "Voting is not active", nil, nil)
	}
	if !now.After(voting.StartDate) && !now.Before(voting.EndDate) {
		return helpers.NewResponse(http.StatusBadRequest, "Voting is not available", nil, nil)
	}

	// check if member already vote
	votingLog, err := u.mongoDbRepo.FetchOneVotingLog(ctx, map[string]interface{}{
		"memberId": claim.UserID,
		"votingId": payload.VotingID,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}
	if votingLog != nil {
		return helpers.NewResponse(http.StatusBadRequest, "Member already vote", nil, nil)
	}

	// create voting log
	newVotingLog := &mongo_model.VotingLog{
		ID: primitive.NewObjectID(),
		Candidate: mongo_model.CandidateFK{
			ID:       candidate.ID.Hex(),
			VotingID: voting.ID.Hex(),
		},
		MemberID:  claim.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = u.mongoDbRepo.CreateOneVotingLog(ctx, newVotingLog)
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// update candidate voter count
	candidate.Voters.Count += 1
	err = u.mongoDbRepo.UpdatePartialCandidate(ctx, map[string]interface{}{
		"id": candidate.ID,
	}, map[string]interface{}{
		"voters.count": candidate.Voters.Count,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	// update voting total vote
	voting.TotalVoter += 1
	err = u.mongoDbRepo.UpdatePartialVoting(ctx, map[string]interface{}{
		"id": voting.ID,
	}, map[string]interface{}{
		"totalVoter": voting.TotalVoter,
	})
	if err != nil {
		return helpers.NewResponse(http.StatusInternalServerError, err.Error(), nil, nil)
	}

	return helpers.NewResponse(http.StatusOK, "Success", nil, nil)
}
