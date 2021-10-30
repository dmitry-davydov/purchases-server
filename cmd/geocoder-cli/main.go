package main

import (
	"flag"
	"github.com/codingsince1985/geo-golang/yandex"
	"github.com/davecgh/go-spew/spew"
	"log"
)

var address string

func init() {
	flag.StringVar(&address, "address", "", "")
}

func main() {
	flag.Parse()

	if len(address) == 0 {
		log.Panic("address is empty")
	}

	g := yandex.Geocoder("faac1e31-94e3-40de-aa3b-82872c4843dd")
	location, err := g.Geocode(address)
	spew.Dump(location)
	spew.Dump(err)
}
