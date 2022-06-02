package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	mydb "github.com/StepanShevelev/parsego/db"
	"log"
	"net/http"
)

func UrlParse() {
	// Request the HTML page.
	res, err := http.Get("https://www.igromania.ru/articles/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".aubli_data a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		url, ok := s.Attr("href")
		if !ok {
			log.Println("error")
		}

		ContentParse(url)

	})

}

//func ParseLink(url string) {
//	resp, err := http.Get(url)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != 200 {
//		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
//	}
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Find the review items
//	doc.Find(".href").Each(func(i int, s *goquery.Selection) {
//		// For each item found, get the title
//		title := s.Find(".aubli_name").Text()
//		fmt.Printf("Review %d: %s\n", i, title)
//	})
//}

func ContentParse(url string) {
	// Request the HTML page.
	res, err := http.Get("https://www.igromania.ru" + url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".page_article_content").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find("div").Text()
		//url, ok := s.Attr("href")
		//if !ok {
		//	log.Println("error")
		//}lcol

		fmt.Printf("Review %d: %s\n", i, title)

	})

}

func main() {
	mydb.ConnectToDb()
	fmt.Println("connected to db")

	UrlParse()
}
