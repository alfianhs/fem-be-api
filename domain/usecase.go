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
}
