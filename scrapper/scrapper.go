package scrapper

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
	url      string
	title    string
	company  string
	location string
	date     string
}

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

func Scrape(term, location string) {
	var baseURL string = "https://www.linkedin.com/jobs/search/?keywords=" + term + "&location=" + location
	var jobs []extractedJob
	extractedJob := getpage(baseURL)
	jobs = append(jobs, extractedJob...)
	//fmt.Println(jobs)
	writeJobs(jobs)
	fmt.Println("writing job is done")

}

func getpage(url string) []extractedJob {
	var jobs []extractedJob

	c := make(chan extractedJob)

	pageURL := url

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
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
	company := card.Find(".result-card__contents>h4").Text()
	location := card.Find(".result-card__meta>.job-result-card__location").Text()
	date := card.Find(".result-card__meta>.job-result-card__listdate").Text()
	//url later
	url, _ := card.Find(".result-card__full-card-link").Attr("href")
	c <- extractedJob{
		id:       id,
		url:      url,
		title:    title,
		company:  company,
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

	headers := []string{"Id", "URL", "Title", "Company", "Location", "Date"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{job.id, job.url, job.title, job.company, job.location, job.date}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}
