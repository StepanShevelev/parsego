package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectToDb() {
	dsn := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable timezone=Europe/Moscow"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
	}

	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&ErrLogs{})

	Database = DbInstance{
		Db: db,
	}
}

//func FindPostByTitle(title string) (*Post, error) {
//	var post *Post
//
//	result := Database.Db.Preload("Pets").Preload("Categories").Find(&post, "title = ?", title)
//
//	if result.Error != nil {
//		UppendErrorWithPath(result.Error)
//		return nil, result.Error
//
//	}
//	return post, nil
//}
//
//func CreatePost(title string, text string, img string) {
//	var post *Post
//
//	post.Title = title
//	post.Text = text
//	//post.Images = img
//	post.Images = append(post.Images, Image{Name: img})
//	//Database.Db.Model(&post).Association("Images").Append(&category)
//	Database.Db.Create(&post)
//	//Database.Db.Model(&category).Association("Users").Append(&user)
//}
