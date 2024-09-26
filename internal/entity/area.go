package entity

type Area struct {
	ParentID int64
	ID       int64  `json:"id"`
	Type     int64  `json:"ateType"`
	TypeName string `json:"nameRu_short"`
	Name     string `json:"ateName"`
}

type AreaResponse struct {
	Children []Area `json:"child"`
}
