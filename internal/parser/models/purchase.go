package models

type Purchase struct {
	Shop       Shop      `json:"shop"`
	Products   []Product `json:"products"`
	TotalPrice int       `json:"totalPrice"`
	Timestamp  int64     `json:"timestamp"`
}
