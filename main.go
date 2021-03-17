package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

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
	var jobs []extractedJob
	extractedJob := getpage()
	jobs = append(jobs, extractedJob...)
	//fmt.Println(jobs)
	writeJobs(jobs)
	fmt.Println("writing job is done")

}

func getpage() []extractedJob {
	var jobs []extractedJob

	c := make(chan extractedJob)

	pageURL := baseURL

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	fmt.Println(doc)
	searchCards := doc.Find(".result-card")
	searchCards.Each(func(i int, card *goquery.Selection) {

		//job := extractJob(card)
		//jobs = append(jobs, job)
		go extractJob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	return jobs
}

//extract job function
//make go routine
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-id")
	title := card.Find(".result-card__contents>h3").Text()
	subtitle := card.Find(".result-card__contents>h4").Text()
	location := card.Find(".result-card__meta>.job-result-card__location").Text()
	date := card.Find(".result-card__meta>.job-result-card__listdate").Text()
	//url later
	c <- extractedJob{
		id:       id,
		title:    title,
		subtitle: subtitle,
		location: location,
		date:     date,
	}
}

//write job function
func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Id", "Title", "Subtitle", "Location", "Date"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.id, job.title, job.subtitle, job.location, job.date}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}
