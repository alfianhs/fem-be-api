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
