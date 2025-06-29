package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"
var pageNumAppear string = "&recruitPage="

type extractedJob struct {
	id          string
	title       string
	location    string
	requirement string
	company     string
}

func main() {
	totalPages := getPages()

	var jobs []extractedJob
	c := make(chan []extractedJob)

	for i := 1; i <= totalPages; i++ {
		go getPage(i, c)
	}

	for i := 1; i <= totalPages; i++ {
		jobs = append(jobs, <-c...)
	}

	writeJob(jobs)
}

// page url 불러오기
func getPage(page int, mainC chan []extractedJob) {
	// page 값을 적용할 URL에 맞게 변경
	pageURL := baseURL + pageNumAppear + strconv.Itoa(page)
	fmt.Println("Request:", pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	var jobs []extractedJob
	c := make(chan extractedJob)

	cards := doc.Find(".item_recruit")
	cards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})

	for i := 0; i < cards.Length(); i++ {
		jobs = append(jobs, <-c)
	}

	mainC <- jobs
}

// source code 에서 pagination 가져와 page 개수 파악
func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), "")
}

func extractJob(s *goquery.Selection, c chan extractedJob) {
	// Attr : attribute(속성), value, id, class, href ...etc
	id, _ := s.Attr("value")
	title := cleanString(s.Find(".job_tit").Text())
	location := cleanString(s.Find(".job_condition>span>a").Text())
	req := s.Find(".job_sector").Clone()
	req.Find(".job_day").Remove()
	requirement := cleanString(req.Text())
	company := cleanString(s.Find(".area_corp>.corp_name").Text())

	c <- extractedJob{
		id:          id,
		title:       title,
		location:    location,
		requirement: requirement,
		company:     company,
	}
}

func writeJob(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Id", "Title", "Location", "Requirement", "Company"}

	Werr := w.Write(headers)
	checkErr(Werr)

	for _, job := range jobs {
		jobSlice := []string{job.id, job.title, job.location, job.requirement, job.company}
		jobErr := w.Write(jobSlice)
		checkErr(jobErr)
	}
}
