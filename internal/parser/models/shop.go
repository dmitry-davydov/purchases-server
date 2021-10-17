package models

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Shop struct {
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}
