package parser

import "github.com/dmitry-davydov/purchases-server/internal/parser/models"

type ParseRequest struct {
	TotalPrice string
	FP string
}

type ParserInterface interface {
	Parse(request ParseRequest) (*models.Purchase, error)
}
