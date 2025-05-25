package mongo_repository

import (
	"app/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDbRepo struct {
	Conn                       *mongo.Database
	superadminCollection       string
	adminCollection            string
	memberCollection           string
	mediaCollection            string
	seasonCollection           string
	venueCollection            string
	teamCollection             string
	playerCollection           string
	seasonTeamCollection       string
	seasonTeamPlayerCollection string
	seriesCollection           string
	ticketCollection           string
	votingCollection           string
	candidateCollection        string
	votingLogCollection        string
}

func NewMongoDbRepo(conn *mongo.Database) domain.MongoDbRepo {
	return &mongoDbRepo{
		Conn:                       conn,
		superadminCollection:       "superadmins",
		adminCollection:            "admins",
		memberCollection:           "members",
		mediaCollection:            "medias",
		seasonCollection:           "seasons",
		venueCollection:            "venues",
		teamCollection:             "teams",
		playerCollection:           "players",
		seasonTeamCollection:       "season_teams",
		seasonTeamPlayerCollection: "season_team_players",
		seriesCollection:           "series",
		ticketCollection:           "tickets",
		votingCollection:           "votings",
		candidateCollection:        "candidates",
		votingLogCollection:        "voting_logs",
	}
}
