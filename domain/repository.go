package domain

import (
	mongo_model "app/domain/model/mongo"
	s3_model "app/domain/model/s3"
	"app/helpers"
	"context"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbRepo interface {
	// Superadmin
	FetchOneSuperadmin(ctx context.Context, options map[string]interface{}) (row *mongo_model.Superadmin, err error)

	// Admin
	FetchOneAdmin(ctx context.Context, options map[string]interface{}) (row *mongo_model.Admin, err error)

	// Member
	FetchOneMember(ctx context.Context, options map[string]interface{}) (row *mongo_model.Member, err error)
	CreateOneMember(ctx context.Context, member *mongo_model.Member) (err error)
	UpdatePartialMember(ctx context.Context, options, field map[string]interface{}) (err error)

	// Media
	FetchOneMedia(ctx context.Context, options map[string]interface{}) (row *mongo_model.Media, err error)
	CreateOneMedia(ctx context.Context, media *mongo_model.Media) (err error)
	CreateManyMedia(ctx context.Context, medias []*mongo_model.Media) (err error)
	UpdatePartialMedia(ctx context.Context, options, field map[string]interface{}) (err error)
	UpdateManyMediaPartial(ctx context.Context, options, field map[string]interface{}) (err error)

	// Season
	FetchListSeason(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountSeason(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneSeason(ctx context.Context, options map[string]interface{}) (row *mongo_model.Season, err error)
	CreateOneSeason(ctx context.Context, season *mongo_model.Season) (err error)
	UpdatePartialSeason(ctx context.Context, options, field map[string]interface{}) (err error)

	// Venue
	FetchListVenue(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountVenue(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneVenue(ctx context.Context, options map[string]interface{}) (row *mongo_model.Venue, err error)
	CreateOneVenue(ctx context.Context, venue *mongo_model.Venue) (err error)
	UpdatePartialVenue(ctx context.Context, options, field map[string]interface{}) (err error)

	// Team
	FetchListTeam(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountTeam(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneTeam(ctx context.Context, options map[string]interface{}) (row *mongo_model.Team, err error)
	CreateOneTeam(ctx context.Context, team *mongo_model.Team) (err error)
	UpdatePartialTeam(ctx context.Context, options, field map[string]interface{}) (err error)

	// Player
	FetchListPlayer(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountPlayer(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOnePlayer(ctx context.Context, options map[string]interface{}) (row *mongo_model.Player, err error)
	CreateOnePlayer(ctx context.Context, venue *mongo_model.Player) (err error)
	UpdatePartialPlayer(ctx context.Context, options, field map[string]interface{}) (err error)

	// Season Team
	FetchListSeasonTeam(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountSeasonTeam(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneSeasonTeam(ctx context.Context, options map[string]interface{}) (row *mongo_model.SeasonTeam, err error)
	CreateManySeasonTeam(ctx context.Context, seasonTeams []*mongo_model.SeasonTeam) (err error)
	UpdatePartialSeasonTeam(ctx context.Context, options, field map[string]interface{}) (err error)
	UpdateManySeasonTeamPartial(ctx context.Context, options, field map[string]interface{}) (err error)

	// Season TeamPlayer
	FetchListSeasonTeamPlayer(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountSeasonTeamPlayer(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneSeasonTeamPlayer(ctx context.Context, options map[string]interface{}) (row *mongo_model.SeasonTeamPlayer, err error)
	CreateOneSeasonTeamPlayer(ctx context.Context, seasonTeamPlayer *mongo_model.SeasonTeamPlayer) (err error)
	UpdatePartialSeasonTeamPlayer(ctx context.Context, options, field map[string]interface{}) (err error)
	UpdateManySeasonTeamPlayerPartial(ctx context.Context, options, field map[string]interface{}) (err error)

	// Series
	FetchListSeries(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountSeries(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneSeries(ctx context.Context, options map[string]interface{}) (row *mongo_model.Series, err error)
	CreateOneSeries(ctx context.Context, season *mongo_model.Series) (err error)
	UpdatePartialSeries(ctx context.Context, options, field map[string]interface{}) (err error)
	SumSeriesMatchCount(ctx context.Context, options map[string]interface{}) (total int64, err error)

	// Ticket
	FetchListTicket(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountTicket(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneTicket(ctx context.Context, options map[string]interface{}) (row *mongo_model.Ticket, err error)
	CreateOneTicket(ctx context.Context, ticket *mongo_model.Ticket) (err error)
	CreateManyTicket(ctx context.Context, tickets []*mongo_model.Ticket) (err error)
	UpdatePartialTicket(ctx context.Context, options, field map[string]interface{}) (err error)
	IncrementOneTicket(ctx context.Context, id string, payload map[string]int64) (err error)

	// Voting
	FetchListVoting(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountVoting(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneVoting(ctx context.Context, options map[string]interface{}) (row *mongo_model.Voting, err error)
	CreateOneVoting(ctx context.Context, voting *mongo_model.Voting) (err error)
	UpdatePartialVoting(ctx context.Context, options, field map[string]interface{}) (err error)

	// Candidate
	FetchListCandidate(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountCandidate(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneCandidate(ctx context.Context, options map[string]interface{}) (row *mongo_model.Candidate, err error)
	CreateOneCandidate(ctx context.Context, candidate *mongo_model.Candidate) (err error)
	UpdatePartialCandidate(ctx context.Context, options, field map[string]interface{}) (err error)

	// Voting Log
	FetchListVotingLog(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountVotingLog(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneVotingLog(ctx context.Context, options map[string]interface{}) (row *mongo_model.VotingLog, err error)
	CreateOneVotingLog(ctx context.Context, candidate *mongo_model.VotingLog) (err error)
	UpdatePartialVotingLog(ctx context.Context, options, field map[string]interface{}) (err error)

	// Purchase
	FetchListPurchase(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountPurchase(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOnePurchase(ctx context.Context, options map[string]interface{}) (row *mongo_model.Purchase, err error)
	CreateOnePurchase(ctx context.Context, purchase *mongo_model.Purchase) (err error)
	UpdatePartialPurchase(ctx context.Context, options, field map[string]interface{}) (err error)

	// Ticket Purchase
	FetchListTicketPurchase(ctx context.Context, options map[string]interface{}) (cur *mongo.Cursor, err error)
	CountTicketPurchase(ctx context.Context, options map[string]interface{}) (total int64)
	FetchOneTicketPurchase(ctx context.Context, options map[string]interface{}) (row *mongo_model.TicketPurchase, err error)
	CreateOneTicketPurchase(ctx context.Context, ticketPurchase *mongo_model.TicketPurchase) (err error)
	CreateManyTicketPurchase(ctx context.Context, ticketPurchases []*mongo_model.TicketPurchase) (err error)
	UpdatePartialTicketPurchase(ctx context.Context, options, field map[string]interface{}) (err error)
	UpdateManyTicketPurchasePartial(ctx context.Context, options, field map[string]interface{}) (err error)
}

type S3Repo interface {
	UploadFilePublic(ctx context.Context, objectName string, body io.Reader, mimeType string) (*s3_model.UploadResponse, error)
}

type XenditRepo interface {
	GenereteSnapLink(ctx context.Context, purchase mongo_model.Purchase) (helpers.Response, error)
}
