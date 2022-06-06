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

const timesleep = 20

func UrlMainParse() {
	// Request the HTML page https://www.igromania.ru/article/32249/Poigrali_v_Destroy_All_Humans!_2-Reprobed_i_delimsya_vpechatleniyami.html
	//https://www.igromania.ru/articles/
	time.Sleep(timesleep * time.Second)
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
		time.Sleep(timesleep * time.Second)
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
			DataParse(*doc, url)

			for _, url := range urlArr {
				ArticleParse(url)
			}

		})

	} else {
		time.Sleep(timesleep * time.Second)
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

			urlArr := FindUrlInArticle(doc)
			DataParse(*doc, url)

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

func DataParse(doc goquery.Document, url string) {

	var image mydb.Image
	var post mydb.Post

	var imgMass [][]byte

	var id = TitleParse(doc)

	doc.Find(".page_article_content").Find(".main_pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		//fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)
		gImg := []byte(img)
		imgMass = append(imgMass, gImg)

		doc.Find(".page_article_content").Each(func(k int, selec *goquery.Selection) {

			if selec.HasClass(".similar_block_body") {
				selec.Remove()
			}
			doc.Find(".pic_container").Find("img").Each(func(j int, selecti *goquery.Selection) {

				img, _ := selecti.Attr("src")
				//fmt.Printf("IMAGE OF ARTICLE %d: %s\n", j, img)

				gImg := []byte(img)
				imgMass = append(imgMass, gImg)
			})
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
		stxt := doc.Find(".cols2").Each(func(i int, article *goquery.Selection) {
			article.Find("ul").Each(func(j int, se *goquery.Selection) {

			})
		}).Text()

		result := mydb.Database.Db.Find(&post, "id = ?", id)
		if result.Error != nil {
			logrus.Info("Error occurred while searching post")
			mydb.UppendErrorWithPath(result.Error)
			return
		}

		result = mydb.Database.Db.Find(&post, "sub_text = ?", stxt)
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info("No such post")
			mydb.UppendErrorWithPath(result.Error)

		}
		post.SubText = stxt

		result = mydb.Database.Db.Save(&post)
		if result.Error != nil {
			logrus.Info("Error occurred while updating post")
			mydb.UppendErrorWithPath(result.Error)
			return
		}
	})

	doc.Find(".lcol ").Each(func(i int, s *goquery.Selection) {
		txt := doc.Find(".page_article_content ").Each(func(i int, article *goquery.Selection) {
			article.Find("span").Each(func(i int, sp *goquery.Selection) {
				sp.Remove()
			})
			article.Find("div").Each(func(j int, se *goquery.Selection) {
				if se.HasClass("similar_block") || se.HasClass("container_wide1") || se.HasClass("uninote") || se.HasClass("plusminus") || se.HasClass("how_achiv") {
					se.Remove()

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
			logrus.Info("No such post")
			mydb.UppendErrorWithPath(result.Error)

		}
		post.Text = txt

		if strings.Contains(url, "https:/") {
			post.ArticleUrl = url
		} else {
			post.ArticleUrl = ("https://www.igromania.ru" + url)

		}

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

		post.Title = title
		result := mydb.Database.Db.Select("Title").Create(&post)
		if result.Error != nil {
			logrus.Info("Error occurred while creating a post")
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
