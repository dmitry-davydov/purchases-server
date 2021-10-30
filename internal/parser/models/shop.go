package models

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Shop struct {
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}
