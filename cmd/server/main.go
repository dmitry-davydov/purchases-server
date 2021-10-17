package main

import (
	"flag"
	"github.com/dmitry-davydov/purchases-server/internal/parser"
	"github.com/labstack/echo/v4"
	"net/http"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "Server Port")
}

func main() {
	flag.Parse()

	srv := parser.NewParserService()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {

		fp := c.QueryParam("fp")
		tp := c.QueryParam("tp")

		if len(fp) == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "fp is empty",
			})
		}

		if len(tp) == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "tp is empty",
			})
		}

		model, err := srv.Parse(fp, tp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
