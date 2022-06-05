package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	mydb "github.com/StepanShevelev/parsego/db"
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

//type PostParse struct {
//	gorm.Model
//	Title  string       `json:"title" db:"title" gorm:"unique"`
//	Text   string       `json:"text" db:"text"`
//	Images []ImageParse `json:"images" db:"images" gorm:"foreignKey:PostID"`
//}

//
//type ImageParse struct {
//
//	Name []byte `json:"name" db:"name"`
//	PostID uint   `json:"post_id" db:"post_id"`
//}

func DataParse(doc *goquery.Document) {

	var image mydb.Image
	var post *mydb.Post
	var id = TitleParse(doc)

	//TitleParse(doc)
	//doc.Find(".page_article").Each(func(i int, s *goquery.Selection) {
	//	title = s.Find(".page_article_ttl").Text()
	//	fmt.Printf("TITLE OF ARTICLE %d: %s\n", i, title)
	//
	//	post.Title = title
	//	mydb.Database.Db.Select("Title").Create(&post)
	//
	//})

	doc.Find(".page_article_content").Find(".main_pic_container").Find("img").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Attr("src")
		//fmt.Printf("IMAGE OF ARTICLE %d: %s\n", i, img)
		gImg := []byte(img)

		image.Name = gImg
		image.PostID = id
		mydb.Database.Db.Create(&image)
		mydb.Database.Db.Find(&post, "id = ?", id)
		logrus.Info(img)
		//
		//mydb.Database.Db.Model(&image).Association("posts").Append(&post)

		doc.Find(".universal_content").Find(".pic_container").Find("img").Each(func(j int, se *goquery.Selection) {
			img, _ := se.Attr("src")
			//fmt.Printf("IMAGE OF ARTICLE %d: %s\n", j, img)

			gImg := []byte(img)
			image.Name = gImg
			image.PostID = id

			mydb.Database.Db.Create(&image)
			mydb.Database.Db.Find(&post, "id = ?", id)
			//logrus.Info(image)

			mydb.Database.Db.Model(&image).Association("posts").Append(&post)

			////Imp.Name = append(Imp.Name, gImg)
			//var images = mydb.Image{{Name: gImg}}
			//mydb.Database.Db.Create(&images)
		})

		//mydb.Database.Db.Model(&category).Association("Users").Append(&user)
		//mydb.Database.Db.Model(&user).Association("Categories").Append(&category)
	})

	doc.Find(".page_article_content ").Each(func(i int, s *goquery.Selection) {
		txt := doc.Find(".page_article_content ").Each(func(i int, article *goquery.Selection) {
			article.Find("div").Each(func(j int, s *goquery.Selection) {
				if s.HasClass("container_wide1") || s.HasClass("uninote") {
					s.Remove()
				}
			})
		}).Text()

		mydb.Database.Db.Find(&post, "id = ?", id)
		//logrus.Info(id)

		post.Text = txt
		mydb.Database.Db.Save(&post)

		//fmt.Printf("TEXT OF ARTICLE : %s\n", txt)
	})

}

func TitleParse(doc *goquery.Document) uint {
	var post mydb.Post

	doc.Find(".page_article").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".page_article_ttl").Text()
		fmt.Printf("TITLE OF ARTICLE %d: %s\n", i, title)

		post.Title = title
		mydb.Database.Db.Select("Title").Create(&post)
		//logrus.Info(post.ID)

		//result := mydb.Database.Db.Find(&post, "id = ?",)
		////logrus.Info(id)
		//if result.Error != nil {
		//	return
		//}
	})

	return post.ID
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

func main() {
	mydb.ConnectToDb()
	fmt.Println("connected to db")

	//UrlMainParse()
	ArticleParse()
	//TitleParse()
	//ImageParse()
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
