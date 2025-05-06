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
	GetSeasonsList(ctx context.Context, options map[string]interface{}) helpers.Response
	GetSeasonDetail(ctx context.Context, options map[string]interface{}) helpers.Response
	CreateSeason(ctx context.Context, payload request.SeasonCreateRequest, request *http.Request) helpers.Response
	UpdateSeason(ctx context.Context, options map[string]interface{}, request *http.Request) helpers.Response
	DeleteSeason(ctx context.Context, options map[string]interface{}) helpers.Response
	UpdateSeasonStatus(ctx context.Context, options map[string]interface{}) helpers.Response

	// Venue
	GetVenueList(ctx context.Context, options map[string]interface{}) helpers.Response
	GetVenueDetail(ctx context.Context, options map[string]interface{}) helpers.Response
	CreateVenue(ctx context.Context, payload request.VenueCreateRequest) helpers.Response
	UpdateVenue(ctx context.Context, options map[string]interface{}) helpers.Response
	DeleteVenue(ctx context.Context, options map[string]interface{}) helpers.Response

	// Team
	GetTeamsList(ctx context.Context, query url.Values) helpers.Response
	GetTeamDetail(ctx context.Context, id string) helpers.Response
	CreateTeam(ctx context.Context, payload request.TeamCreateRequest, request *http.Request) helpers.Response
	UpdateTeam(ctx context.Context, id string, payload request.TeamUpdateRequest, request *http.Request) helpers.Response
	DeleteTeam(ctx context.Context, id string) helpers.Response
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
