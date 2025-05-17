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
	SeasonStatusDraft     SeriesStatus = 1
	SeriesStatusActive    SeriesStatus = 2
	SeriesStatusNonActive SeriesStatus = 3
)

type SeriesStatusStruct struct {
	ID   SeriesStatus `json:"id"`
	Name string       `json:"name"`
}

var SeriesStatusMap = map[SeriesStatus]SeriesStatusStruct{
	SeasonStatusDraft:     {ID: SeasonStatusDraft, Name: "Draft"},
	SeriesStatusActive:    {ID: SeriesStatusActive, Name: "Active"},
	SeriesStatusNonActive: {ID: SeriesStatusNonActive, Name: "Non-Active"},
}
