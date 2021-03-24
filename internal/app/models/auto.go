package models

//Auto models...
type Auto struct {
	ID       int    `json:"id"`
	Mark     string `json:"mark"`
	MaxSpeed int64  `json:"max_speed"`
	Distance int64  `json:"distance"`
	Handler  string `json:"handler"`
	Stock    string `json:"stock"`
}

type AutoApi struct {
	MaxSpeed int64  `json:"max_speed"`
	Distance int64  `json:"distance"`
	Handler  string `json:"handler"`
	Stock    string `json:"stock"`
}
