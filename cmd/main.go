package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

func UrlMainParse() {
	// Request the HTML page
	time.Sleep(5 * time.Second)
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

	// Find the review items Find(".aubli_data a")
	doc.Find(".aubli_data a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		Url, ok := s.Attr("href")
		if !ok {
			log.Println("error")
		}
		ArticleParse(Url)
		fmt.Printf("ARTICLE URL %d: %s\n", i, Url)

	})
}

func ArticleParse(Url string) {
	// Request the HTML page
	if strings.Contains(Url, "https:/") {
		time.Sleep(5 * time.Second)
		res, err := http.Get(Url)
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

		doc.Find(".page_article_content a").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the title
			TextParse(doc)
			TitleParse(doc)
			ImageParse(doc)
			FindUrlInArticle(doc)

		})

	} else {
		time.Sleep(5 * time.Second)
		res, err := http.Get("https://www.igromania.ru" + Url)
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
		doc.Find(".page_article_content a").Each(func(i int, s *goquery.Selection) {

			// For each item found, get the title
			DataParse(doc)
			FindUrlInArticle(doc)

			//TextParse(doc)
			//ImageParse(doc)

		})
	}

}

func DataParse(doc *goquery.Document) {

	doc.Find(".page_article").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find(".page_article_ttl").Text()

		fmt.Printf("TITLE OF ARTICLE %d: %s\n", i, title)

	})

	doc.Find(".page_article_content").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		text := s.Find("div").Text()

		fmt.Printf("TEXT OF ARTICLE %d: %s\n", i, text)
	})

	doc.Find(".universal_content").Find(".pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title

		img, _ := s.Attr("src")

		fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)

	})

}

//func TextParse(doc *goquery.Document) {
//
//	// Find the review items
//	doc.Find(".page_article_content").Each(func(i int, s *goquery.Selection) {
//		// For each item found, get the title
//		text := s.Find("div").Text()
//
//		fmt.Printf("TEXT OF ARTICLE %d: %s\n", i, text)
//	})
//
//}
//
//func TitleParse(doc *goquery.Document) {
//
//	// Find the review items
//	doc.Find(".page_article").Each(func(i int, s *goquery.Selection) {
//		// For each item found, get the title
//		title := s.Find(".page_article_ttl").Text()
//
//		fmt.Printf("TITLE OF ARTICLE %d: %s\n", i, title)
//
//	})
//}
//
//func ImageParse(doc *goquery.Document) {
//	// Find the review items  pic_container
//	doc.Find(".universal_content").Find(".pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
//		// For each item found, get the title
//
//		img, _ := s.Attr("src")
//
//		fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)
//
//	})
//}

func FindUrlInArticle(doc *goquery.Document) {

	// Find the review items Find(".aubli_data a")
	doc.Find(".uninote console a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		Url, ok := s.Attr("href")
		if !ok {
			log.Println("error")
		}
		//ArticleParse(Url)
		fmt.Printf("NEW LINK FROM ARTICLE %d: %s\n", i, Url)

	})
}

func main() {
	//mydb.ConnectToDb()
	//fmt.Println("connected to db")

	//ArticleUrlParse()
	UrlMainParse()
	//TitleParse()
	//ImageParse()
}
