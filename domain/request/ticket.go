package request

type TicketCreateOrUpdateRequest struct {
	SeriesID string          `json:"seriesId"`
	Tickets  []TicketRequest `json:"tickets"`
}

type TicketRequest struct {
	ID     string               `json:"id"`
	Name   string               `json:"name"`
	Date   string               `json:"date"`
	Price  float64              `json:"price"`
	Quota  int64                `json:"quota"`
	Matchs []TicketMatchRequest `json:"matchs"`
}

type TicketMatchRequest struct {
	HomeSeasonTeamID string `json:"homeSeasonTeamId"`
	AwaySeasonTeamID string `json:"awaySeasonTeamId"`
	Time             string `json:"time"`
}
