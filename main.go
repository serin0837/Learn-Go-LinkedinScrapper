package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type extractedReveiw struct {
	id       string
	location string
	title    string
}

var baseURL string = "https://www.linkedin.com/jobs/search/?keywords=python"

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
	totalPages := getPages()
	fmt.Println(totalPages)
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	fmt.Println(doc)
	doc.Find(".artdeco-pagination").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Find("artdeco-pagination__indicator"))
	})
	return pages
}
