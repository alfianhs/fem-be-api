package domain

import (
	"app/domain/request"
	"app/helpers"
	jwt_helpers "app/helpers/jwt"
	"context"
	"net/http"
	"net/url"
)

type SuperadminAppUsecase interface {
	// Auth
	Login(ctx context.Context, payload request.SuperadminLoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim jwt_helpers.SuperadminJWTClaims) helpers.Response

	// Season
	GetSeasonsList(ctx context.Context, query url.Values) helpers.Response
	GetSeasonDetail(ctx context.Context, id string) helpers.Response
	GetActiveSeasonDetail(ctx context.Context) helpers.Response
	CreateSeason(ctx context.Context, payload request.SeasonCreateRequest, request *http.Request) helpers.Response
	UpdateSeason(ctx context.Context, id string, payload request.SeasonUpdateRequest, request *http.Request) helpers.Response
	DeleteSeason(ctx context.Context, id string) helpers.Response
	UpdateSeasonStatus(ctx context.Context, id string, payload request.SeasonStatusUpdateRequest) helpers.Response

	// Venue
	GetVenueList(ctx context.Context, query url.Values) helpers.Response
	GetVenueDetail(ctx context.Context, id string) helpers.Response
	CreateVenue(ctx context.Context, payload request.VenueCreateRequest) helpers.Response
	UpdateVenue(ctx context.Context, id string, payload request.VenueUpdateRequest) helpers.Response
	DeleteVenue(ctx context.Context, id string) helpers.Response

	// Team
	GetTeamsList(ctx context.Context, query url.Values) helpers.Response
	GetTeamDetail(ctx context.Context, id string) helpers.Response
	CreateTeam(ctx context.Context, payload request.TeamCreateRequest, request *http.Request) helpers.Response
	UpdateTeam(ctx context.Context, id string, payload request.TeamUpdateRequest, request *http.Request) helpers.Response
	DeleteTeam(ctx context.Context, id string) helpers.Response

	// Player
	GetPlayerList(ctx context.Context, query url.Values) helpers.Response
	GetPlayerDetail(ctx context.Context, id string) helpers.Response
	CreatePlayer(ctx context.Context, payload request.PlayerCreateRequest) helpers.Response
	UpdatePlayer(ctx context.Context, id string, payload request.PlayerUpdateRequest) helpers.Response
	DeletePlayer(ctx context.Context, id string) helpers.Response

	// SeasonTeam
	GetSeasonTeamsList(ctx context.Context, query url.Values) helpers.Response
	GetSeasonTeamDetail(ctx context.Context, id string) helpers.Response
	CreateSeasonTeam(ctx context.Context, payload request.SeasonTeamCreateRequest) helpers.Response
	DeleteSeasonTeam(ctx context.Context, id string) helpers.Response
	ManageSeasonTeam(ctx context.Context, payload request.SeasonTeamManageRequest) helpers.Response

	// SeasonTeamPlayer
	GetSeasonTeamPlayersList(ctx context.Context, query url.Values) helpers.Response
	GetPlayerPositionsList(ctx context.Context) helpers.Response
	GetSeasonTeamPlayerDetail(ctx context.Context, id string) helpers.Response
	CreateSeasonTeamPlayer(ctx context.Context, payload request.SeasonTeamPlayerCreateRequest, request *http.Request) helpers.Response
	UpdateSeasonTeamPlayer(ctx context.Context, id string, payload request.SeasonTeamPlayerUpdateRequest, request *http.Request) helpers.Response
	DeleteSeasonTeamPlayer(ctx context.Context, id string) helpers.Response

	// Series
	GetSeriesList(ctx context.Context, queryParam url.Values) helpers.Response
	GetSeriesDetail(ctx context.Context, id string) helpers.Response
	CreateSeries(ctx context.Context, payload request.SeriesCreateRequest) helpers.Response
	UpdateSeries(ctx context.Context, id string, payload request.SeriesUpdateRequest) helpers.Response
	DeleteSeries(ctx context.Context, id string) helpers.Response

	// Ticket
	GetTicketsList(ctx context.Context, queryParam url.Values) helpers.Response
	GetTicketDetail(ctx context.Context, id string) helpers.Response
	CreateOrUpdateTicket(ctx context.Context, payload request.TicketCreateOrUpdateRequest) helpers.Response
	DeleteTicket(ctx context.Context, id string) helpers.Response

	// Voting
	GetVotingList(ctx context.Context, queryParam url.Values) helpers.Response
	GetVotingDetail(ctx context.Context, id string) helpers.Response
	CreateVoting(ctx context.Context, payload request.VotingCreateRequest, request *http.Request) helpers.Response
	UpdateVoting(ctx context.Context, id string, payload request.VotingUpdateRequest, request *http.Request) helpers.Response
	DeleteVoting(ctx context.Context, id string) helpers.Response

	// Candidate

	GetCandidateList(ctx context.Context, q url.Values) helpers.Response
	GetCandidateDetail(ctx context.Context, id string) helpers.Response
	CreateCandidate(ctx context.Context, req request.CandidateCreateRequest) helpers.Response
	UpdateCandidate(ctx context.Context, id string, req request.CandidateUpdateRequest) helpers.Response
	DeleteCandidate(ctx context.Context, id string) helpers.Response
}

type AdminAppUsecase interface {
	// Auth
	Login(ctx context.Context, payload request.AdminLoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim jwt_helpers.AdminJWTClaims) helpers.Response
}

type MemberAppUsecase interface {
	// Auth
	Register(ctx context.Context, payload request.MemberRegisterRequest) helpers.Response
	VerifyEmail(ctx context.Context, payload request.VerifyEmailRequest) helpers.Response
	ResendEmailVerification(ctx context.Context, payload request.ResendEmailVerificationRequest) helpers.Response
	Login(ctx context.Context, payload request.MemberLoginRequest) helpers.Response
	GetProfile(ctx context.Context, claim jwt_helpers.MemberJWTClaims) helpers.Response

	// Voting
	GetVotingList(ctx context.Context, queryParam url.Values) helpers.Response

	// Candidate
	GetCandidateList(ctx context.Context, claim jwt_helpers.MemberJWTClaims, q url.Values) helpers.Response
	CandidateVote(ctx context.Context, claim jwt_helpers.MemberJWTClaims, payload request.CandidateVoteRequest) helpers.Response

	// Purchase
	CreatePurchase(ctx context.Context, claim jwt_helpers.MemberJWTClaims, payload request.CreatePurchaseRequest) helpers.Response
}
