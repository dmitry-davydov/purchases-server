package main

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/dmitry-davydov/purchases-server/internal/parser"
	"log"
)

var fp string
var totalPrice string

func init() {
	flag.StringVar(&fp, "fp", "", "--fp=1499427403")
	flag.StringVar(&totalPrice, "tp", "", "--tp=492")
}

func main() {
	flag.Parse()

	if len(fp) == 0 {
		log.Panic("fp is empty")
	}

	if len(totalPrice) == 0 {
		log.Panic("totalPrice is empty")
	}

	service := parser.NewParserService()
	purchase, err := service.Parse(fp, totalPrice)

	if err != nil {
		log.Panic(err.Error())
	}

	spew.Dump(purchase)
}
