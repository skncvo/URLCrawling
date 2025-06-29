package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

type extractedJob struct {
	id          string
	title       string
	location    string
	requirement string
	company     string
}

func main() {
	totalPages := getPages()
	for i := 1; i <= totalPages; i++ {
		getPage(i)
		fmt.Println()
	}

}

// page url 불러오기
func getPage(page int) {
	pageURL := baseURL + "&recruitPage=" + strconv.Itoa(page)
	fmt.Println("Request:", pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// https://www.saramin.co.kr/zf_user/jobs/relay/view?rec_idx=공고ID

	doc.Find(".item_recruit").Each(func(i int, s *goquery.Selection) {
		extractJob(s)
	})
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

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), "")
}

func extractJob(s *goquery.Selection) {
	// Attr : attribute(속성), value, id, class, href ...etc
	id, _ := s.Attr("value")
	title := cleanString(s.Find(".job_tit").Text())
	location := cleanString(s.Find(".job_condition>span>a").Text())
	req := s.Find(".job_sector").Clone()
	req.Find(".job_day").Remove()
	requirement := cleanString(req.Text())
	company := cleanString(s.Find(".area_corp>.corp_name").Text())

	fmt.Println(id, title, location, requirement, company)
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
