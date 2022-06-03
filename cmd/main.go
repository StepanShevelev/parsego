package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

func UrlMainParse() {
	// Request the HTML page https://www.igromania.ru/article/32249/Poigrali_v_Destroy_All_Humans!_2-Reprobed_i_delimsya_vpechatleniyami.html
	//https://www.igromania.ru/articles/
	time.Sleep(5 * time.Second)
	res, err := http.Get("https://www.igromania.ru/article/32249/Poigrali_v_Destroy_All_Humans!_2-Reprobed_i_delimsya_vpechatleniyami.html")
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
		//ArticleParse(Url)
		fmt.Printf("ARTICLE URL %d: %s\n", i, Url)

	})
}

// Url string

func ArticleParse() {
	// Request the HTML page
	//if strings.Contains(Url, "https:/") {
	//	time.Sleep(5 * time.Second)
	//	res, err := http.Get(Url)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer res.Body.Close()
	//	if res.StatusCode != 200 {
	//		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	//	}
	//	// Load the HTML document
	//	doc, err := goquery.NewDocumentFromReader(res.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	doc.Find(".page_article_content a").Each(func(i int, s *goquery.Selection) {
	//		// For each item found, get the title
	//		DataParse(doc)
	//		FindUrlInArticle(doc)
	//
	//	})
	//
	//} else {
	//	time.Sleep(5 * time.Second)
	//	res, err := http.Get("https://www.igromania.ru" + Url)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer res.Body.Close()
	//	if res.StatusCode != 200 {
	//		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	//	}
	//
	//	// Load the HTML document
	//	doc, err := goquery.NewDocumentFromReader(res.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	doc.Find(".page_article_content a").Each(func(i int, s *goquery.Selection) {
	//
	//		// For each item found, get the title
	//		DataParse(doc)
	//		FindUrlInArticle(doc)
	//
	//	})
	//}

	time.Sleep(5 * time.Second)
	res, err := http.Get("https://www.igromania.ru/article/32249/Poigrali_v_Destroy_All_Humans!_2-Reprobed_i_delimsya_vpechatleniyami.html")
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
	DataParse(doc)
	FindUrlInArticle(doc)

}

func DataParse(doc *goquery.Document) {

	doc.Find(".page_article").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".page_article_ttl").Text()
		fmt.Printf("TITLE OF ARTICLE %d: %s\n", i, title)
	})

	doc.Find(".page_article_content").Find(".main_pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)
	})

	doc.Find(".page_article_content").Find(".pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)
	})

	doc.Find(".page_article_content ").Each(func(i int, s *goquery.Selection) {
		txt := doc.Find(".page_article_content ").Each(func(i int, article *goquery.Selection) {
			article.Find("div").Each(func(j int, s *goquery.Selection) {
				if s.HasClass("container_wide1") || s.HasClass("uninote") {
					s.Remove()
				}
			})
		}).Text()

		fmt.Printf("TEXT OF ARTICLE : %s\n", txt)
	})

}

func FindUrlInArticle(doc *goquery.Document) {

	logrus.Info("FindUrlInArticle starts")

	doc.Find(".uninote a").Each(func(i int, s *goquery.Selection) {

		Url, ok := s.Attr("href")
		if !ok {
			log.Println("error")
		}
		//ArticleParse(Url)
		fmt.Printf("NEW LINK FROM ARTICLE %d: %s\n", i, Url)

	})

	logrus.Info("FindUrlInArticle finishes")
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

func main() {
	//mydb.ConnectToDb()
	//fmt.Println("connected to db")

	//UrlMainParse()
	ArticleParse()
	//TitleParse()
	//ImageParse()
}
