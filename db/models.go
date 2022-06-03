package db

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title  string  `json:"title" db:"title" gorm:"unique"`
	Text   string  `json:"text" db:"text"`
	Images []Image `json:"pets" db:"pets" gorm:"foreignKey:PostID"`
}

type Image struct {
	gorm.Model
	Name   string `json:"name" db:"name"`
	PostID int    `json:"post_id" db:"post_id"`
}

// ErrLogs storage some error logs
type ErrLogs struct {
	gorm.Model
	Error string `json:"error" db:"error"`
	Place string `json:"place" db:"place"`
	Count int    `json:"count" db:"count"`
}
