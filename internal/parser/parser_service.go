package parser

import (
	"fmt"
	"github.com/dmitry-davydov/purchases-server/internal/parser/models"
	"strings"
	"sync"
)

type response struct {
	model *models.Purchase
	err error
}

type ParserService struct {
	parsers []ParserInterface
}

func (t *ParserService) Parse(FP string, totalPrice string) (*models.Purchase, error) {

	var wg sync.WaitGroup
	var responses []response
	var mut sync.Mutex

	parseRequest := ParseRequest{
		TotalPrice: totalPrice,
		FP:         FP,
	}

	for _, parser := range t.parsers {
		wg.Add(1)

		go func(wg *sync.WaitGroup, parser ParserInterface) {
			model, err := parser.Parse(parseRequest)
			mut.Lock()

			responses = append(responses, response {
				model: model,
				err: err,
			})
			mut.Unlock()
			wg.Done()
		}(&wg, parser)
	}

	wg.Wait()
	var responseErrors []string

	// обработать респонсы
	for _, response := range responses {
		if response.err == nil {
			return response.model, nil
		}

		responseErrors = append(responseErrors, response.err.Error())
	}

	return nil, fmt.Errorf("Can not parse: %s", strings.Join(responseErrors, ", "))
}

func NewParserService() *ParserService {
	t := new(ParserService)
	t.parsers = []ParserInterface{NewTaxcom()}

	return t
}