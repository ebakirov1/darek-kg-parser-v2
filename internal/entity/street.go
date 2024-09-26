package entity

type Street struct {
	AreaID   int64
	StreetID int64  `json:"id"`
	Type     int64  `json:"type"`
	Name     string `json:"name"`
	NameTP   string `json:"nameTp"`
}
