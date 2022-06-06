package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	mydb "github.com/StepanShevelev/parsego/db"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

func UrlMainParse() {
	// Request the HTML page https://www.igromania.ru/article/32249/Poigrali_v_Destroy_All_Humans!_2-Reprobed_i_delimsya_vpechatleniyami.html
	//https://www.igromania.ru/articles/
	time.Sleep(20 * time.Second)
	res, err := http.Get("https://www.igromania.ru/articles/")
	if err != nil {
		logrus.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {

		logrus.Fatalf("status code error: %d %s", res.StatusCode, res.Status)

	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	// Find the review items Find(".aubli_data a")
	doc.Find(".aubli_data a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		Url, ok := s.Attr("href")
		if !ok {
			logrus.Info("error, articles not found")
		}
		ArticleParse(Url)
		fmt.Printf("ARTICLE URL %d: %s\n", i, Url)

	})
}

// Url string

func ArticleParse(url string) {

	if strings.Contains(url, "https:/") {
		time.Sleep(20 * time.Second)
		res, err := http.Get(url)
		if err != nil {
			logrus.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			logrus.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".page_article_content a").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the title
			urlArr := FindUrlInArticle(doc)
			DataParse(*doc)

			for _, url := range urlArr {
				ArticleParse(url)
			}

		})

	} else {
		time.Sleep(20 * time.Second)
		res, err := http.Get("https://www.igromania.ru" + url)
		if err != nil {
			logrus.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			logrus.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			logrus.Fatal(err)
		}
		doc.Find(".page_article_content a").Each(func(i int, s *goquery.Selection) {

			// For each item found, get the title
			urlArr := FindUrlInArticle(doc)
			DataParse(*doc)

			for _, url := range urlArr {
				ArticleParse(url)
			}

		})
	}

	//time.Sleep(5 * time.Second)
	//res, err := http.Get("https://www.igromania.ru/article/32249/Poigrali_v_Destroy_All_Humans!_2-Reprobed_i_delimsya_vpechatleniyami.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer res.Body.Close()
	//if res.StatusCode != 200 {
	//	log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	//}
	//
	//// Load the HTML document
	//doc, err := goquery.NewDocumentFromReader(res.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//urlArr := FindUrlInArticle(doc)
	//DataParse(*doc)
	//
	//for _, url := range urlArr {
	//	ArticleParse(url)
	//}

}

func DataParse(doc goquery.Document) {

	var image mydb.Image
	var post mydb.Post

	var imgMass [][]byte

	var id = TitleParse(doc)

	doc.Find(".page_article_content").Find(".main_pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		//fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)
		gImg := []byte(img)
		imgMass = append(imgMass, gImg)

		doc.Find(".universal_content").Find(".pic_container").Find("img").Each(func(j int, se *goquery.Selection) {
			img, _ := se.Attr("src")
			//fmt.Printf("IMAGE OF ARTICLE %d: %s\n", j, img)

			gImg := []byte(img)
			imgMass = append(imgMass, gImg)
		})
	})

	for _, gImg := range imgMass {
		image.Name = gImg
		image.PostID = id
		var images = []mydb.Image{{Name: gImg, PostID: id}}
		result := mydb.Database.Db.Create(&images)
		if result.Error != nil {
			logrus.Info("Error occurred while creating an image")
			mydb.UppendErrorWithPath(result.Error)
			return
		}
		result = mydb.Database.Db.Find(&post, "id = ?", id)
		if result.Error != nil {
			logrus.Info("Error occurred while searching post")
			mydb.UppendErrorWithPath(result.Error)
			return
		}

	}

	doc.Find(".page_article_content ").Each(func(i int, s *goquery.Selection) {
		txt := doc.Find(".page_article_content ").Each(func(i int, article *goquery.Selection) {
			article.Find("div").Each(func(j int, s *goquery.Selection) {
				if s.HasClass("container_wide1") || s.HasClass("uninote") {
					s.Remove()
				}
			})
		}).Text()

		result := mydb.Database.Db.Find(&post, "id = ?", id)
		if result.Error != nil {
			logrus.Info("Error occurred while searching post")
			mydb.UppendErrorWithPath(result.Error)
			return
		}

		result = mydb.Database.Db.Find(&post, "text = ?", txt)
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info("Post already exist ")
			mydb.UppendErrorWithPath(result.Error)
			return
		}
		post.Text = txt
		result = mydb.Database.Db.Save(&post)
		if result.Error != nil {
			logrus.Info("Error occurred while updating post")
			mydb.UppendErrorWithPath(result.Error)
			return
		}

	})

}

func TitleParse(doc goquery.Document) uint {
	var post mydb.Post

	doc.Find(".page_article").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".page_article_ttl").Text()
		fmt.Printf("TITLE OF ARTICLE %d: %s\n", i, title)

		post.Title = title
		result := mydb.Database.Db.Select("Title").Create(&post)
		if result.Error != nil {
			logrus.Info("Error occurred while creating an image")
			mydb.UppendErrorWithPath(result.Error)
			return
		}

	})

	return post.ID
}

func FindUrlInArticle(doc *goquery.Document) []string {
	logrus.Info("FindUrlInArticle starts")

	var urlArr []string

	doc.Find(".uninote a").Each(func(i int, s *goquery.Selection) {

		Url, ok := s.Attr("href")
		if !ok {
			logrus.Info("could not find related articles")
		}

		urlArr = append(urlArr, Url)
		fmt.Printf("NEW LINK FROM ARTICLE %d: %s\n", i, Url)

	})

	logrus.Info("FindUrlInArticle finishes")
	return urlArr
}

func main() {
	mydb.ConnectToDb()
	logrus.Info("Connected to db")
	UrlMainParse()
}
