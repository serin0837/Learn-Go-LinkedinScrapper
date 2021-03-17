package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	subtitle string
	location string
	date     string
}

//in this page there is no pagination when I am not login
// add location and keyword function
var baseURL string = "https://www.linkedin.com/jobs/search/?keywords=python&location=United%20Kingdom"

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("request failed with Status :", res.StatusCode)
	}
}

func main() {
	getpage()
}

func getpage() {
	pageURL := baseURL

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	fmt.Println(doc)
	searchCards := doc.Find(".result-card")
	searchCards.Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("data-id")
		fmt.Println(id, "<-id")
		title := s.Find(".result-card__contents>h3").Text()
		fmt.Println(title, "<-title")
		subtitle := s.Find(".result-card__contents>h4").Text()
		fmt.Println(subtitle, "<-subtitle")
		location := s.Find(".result-card__meta>.job-result-card__location").Text()
		fmt.Println(location, "<-location")
		date := s.Find(".result-card__meta>.job-result-card__listdate").Text()
		fmt.Println(date, "<-date")
		//url later
	})
}
