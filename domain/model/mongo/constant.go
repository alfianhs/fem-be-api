package mongo_model

type SeasonStatus int

const (
	SeasonStatusActive   SeasonStatus = 1
	SeasonStatusInactive SeasonStatus = 2
)

type SeasonStatusStruct struct {
	ID   SeasonStatus `json:"id"`
	Name string       `json:"name"`
}

var SeasonStatusMap = map[SeasonStatus]SeasonStatusStruct{
	SeasonStatusActive:   {ID: SeasonStatusActive, Name: "Active"},
	SeasonStatusInactive: {ID: SeasonStatusInactive, Name: "Inactive"},
}

var PlayerPositionList = []string{
	"Goalkeeper",
	"Flank",
	"Pivot",
	"Anchor",
}

type SeriesStatus int

const (
	SeriesStatusActive    SeriesStatus = 1
	SeriesStatusNonActive SeriesStatus = 2
	SeriesStatusDraft     SeriesStatus = 3
)

type SeriesStatusStruct struct {
	ID   SeriesStatus `json:"id"`
	Name string       `json:"name"`
}

var SeriesStatusMap = map[SeriesStatus]SeriesStatusStruct{
	SeriesStatusActive:    {ID: SeriesStatusActive, Name: "Active"},
	SeriesStatusNonActive: {ID: SeriesStatusNonActive, Name: "Non-Active"},
	SeriesStatusDraft:     {ID: SeriesStatusDraft, Name: "Draft"},
}

type VotingStatus int

const (
	VotingStatusActive     VotingStatus = 1
	VotingStatusNonActive  VotingStatus = 2
	VotingStatusComingSoon VotingStatus = 3
)

type VotingStatusStruct struct {
	ID   VotingStatus `json:"id"`
	Name string       `json:"name"`
}

var VotingStatusMap = map[VotingStatus]VotingStatusStruct{
	VotingStatusActive:     {ID: VotingStatusActive, Name: "Active"},
	VotingStatusNonActive:  {ID: VotingStatusNonActive, Name: "Non-Active"},
	VotingStatusComingSoon: {ID: VotingStatusComingSoon, Name: "Coming Soon"},
}

type PurchaseStatus int

const (
	PurchaseStatusPending   PurchaseStatus = 1
	PurchaseStatusCompleted PurchaseStatus = 2
	PurchaseStatusFailed    PurchaseStatus = 3
)

type PurchaseStatusStruct struct {
	ID   PurchaseStatus `json:"id"`
	Name string         `json:"name"`
}

var PurchaseStatusMap = map[PurchaseStatus]PurchaseStatusStruct{
	PurchaseStatusPending:   {ID: PurchaseStatusPending, Name: "Pending"},
	PurchaseStatusCompleted: {ID: PurchaseStatusCompleted, Name: "Completed"},
	PurchaseStatusFailed:    {ID: PurchaseStatusFailed, Name: "Failed"},
}
