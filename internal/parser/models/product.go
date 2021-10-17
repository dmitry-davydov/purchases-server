package models

type Product struct {
	Name            string `json:"name"`
	Quantity        int `json:"quantity"`
	PriceForOneItem int `json:"oneItemPrice"`
	TotalPrice      int `json:"totalPrice"`
	VAT             string `json:"vat"`
}
