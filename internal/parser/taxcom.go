package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/yandex"
	"github.com/dmitry-davydov/purchases-server/internal/parser/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Taxcom struct {
	geocoder geo.Geocoder
}

func (t *Taxcom) Parse(request ParseRequest) (*models.Purchase, error) {
	resp, err := http.Get(fmt.Sprintf("https://receipt.taxcom.ru/v01/show?fp=%s&s=%s&sf=False&sfn=False", request.FP, request.TotalPrice))
	if err != nil {

		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	shopAddress := ""
	if shopAddressSelection := doc.Find(".receipt-header2.receipt-company-name .value.receipt-value-1009").First(); shopAddressSelection != nil {
		shopAddress = shopAddressSelection.Text()
	}

	shopName := ""
	if shopNameSelection := doc.Find(".receipt-header1 .receipt-subtitle a").First(); shopNameSelection != nil {
		shopName = shopNameSelection.Text()
	}

	shop := models.Shop{
		Name:    shopName,
		Address: shopAddress,
		Coordinates: models.Coordinates{
			Latitude:  0,
			Longitude: 0,
		},
	}

	purchase := models.Purchase{
		Shop:       shop,
		Products:   []models.Product{},
		TotalPrice: 0,
		Timestamp:  0,
	}

	doc.Find(".receipt-body .items .item").Each(func(i int, selection *goquery.Selection) {

		productName := ""
		if productNameSelection := selection.Find(".value.receipt-value-1030").First(); productNameSelection != nil {
			productName = strings.TrimSpace(productNameSelection.Text())
		}

		quantity := 0
		if quantitySelection := selection.Find(".value.receipt-value-1023").First(); quantitySelection != nil {
			if intVar, err := strconv.Atoi(quantitySelection.Text()); err == nil {
				quantity = intVar
			}
		}

		priceForOne := 0
		if priceForOneSelection := selection.Find(".value.receipt-value-1079").First(); priceForOneSelection != nil {
			if price, err := strconv.ParseFloat(priceForOneSelection.Text(), 64); err == nil {
				priceForOne = int(price * 100)
			}
		}

		vat := ""
		if vatSelection := selection.Find(".name.receipt-value-1199").First(); vatSelection != nil {
			vat = vatSelection.Text()
		}

		product := models.Product{
			Name:            productName,
			Quantity:        quantity,
			PriceForOneItem: priceForOne,
			TotalPrice:      priceForOne * quantity,
			VAT:             vat,
			Category:        "Продукты",
		}

		purchase.Products = append(purchase.Products, product)
	})

	totalPriceNumbers := strings.Trim(request.TotalPrice, ".")
	if totalPriceFloat, err := strconv.ParseFloat(totalPriceNumbers, 64); err == nil {
		purchase.TotalPrice = int(totalPriceFloat * 100)
	}

	if purchaseTimeSelection := doc.Find(".receipt-header2 .value.receipt-value-1012").First(); purchaseTimeSelection != nil {
		timeParts := strings.Split(purchaseTimeSelection.Text(), " ")

		dParts := strings.Split(timeParts[0], ".")
		tParts := strings.Split(timeParts[1], ":")
		dInt := 0
		mInt := 0
		yInt := 0

		hInt := 0
		mmInt := 0

		if dd, err := strconv.Atoi(dParts[0]); err == nil {
			dInt = dd
		}
		if mm, err := strconv.Atoi(dParts[1]); err == nil {
			mInt = mm
		}
		if yy, err := strconv.Atoi("20" + dParts[2]); err == nil {
			yInt = yy
		}
		if hh, err := strconv.Atoi(tParts[0]); err == nil {
			hInt = hh
		}
		if m, err := strconv.Atoi(tParts[1]); err == nil {
			mmInt = m
		}

		date := time.Date(yInt, time.Month(mInt), dInt, hInt, mmInt, 0, 0, time.UTC)
		purchase.Timestamp = date.Unix()
	}

	if location, err := t.geocoder.Geocode(purchase.Shop.Address); err != nil {
		log.Print(err.Error())
	} else {
		purchase.Shop.Coordinates.Longitude = location.Lng
		purchase.Shop.Coordinates.Latitude = location.Lat
	}

	return &purchase, nil
}

func NewTaxcom() *Taxcom {
	t := new(Taxcom)
	t.geocoder = yandex.Geocoder("faac1e31-94e3-40de-aa3b-82872c4843dd")
	return t
}
