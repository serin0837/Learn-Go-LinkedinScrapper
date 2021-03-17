package main

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/serin0837/learngo-linkedin/scrapper"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove("jobs.csv")
	term := strings.ToLower(c.FormValue("term"))
	location := strings.ToLower(c.FormValue("location"))
	scrapper.Scrape(term, location)
	return c.Attachment("jobs.csv", "jobs.csv")
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":4000"))
}
