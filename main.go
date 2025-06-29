package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

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
		// Attr : attribute(속성), value, id, class, href ...etc
		id, _ := s.Attr("value")

		// job_tit의 a를 찾는다
		title := s.Find(".job_tit>a").Text()
		condition := s.Find(".job_condition").Text()

		// 특정 span 적용 안되게
		jobsel := s.Find(".job_sector").Clone()
		jobsel.Find(".job_day").Remove()
		job_sector := jobsel.Text()
		fmt.Println(id, title, condition, job_sector)

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
