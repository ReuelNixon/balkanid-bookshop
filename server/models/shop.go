package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ID 			   	   uint   `gorm:"primaryKey"`
	ISBN               string `json:"isbn" gorm:"not null"`
	BookTitle          string `json:"book_title" gorm:"not null"`
	BookAuthor         string `json:"book_author"`
	YearOfPublication  string `json:"year_of_publication"`
	Publisher          string `json:"publisher"`
	ImageURL           string `json:"image_url"`
}

type Inventory struct {
	gorm.Model
	ID 			   	   uint   `gorm:"primaryKey"`
	BookID			   uint   `json:"book_id" gorm:"not null"`
	Quantity		   uint   `json:"quantity" gorm:"not null"`
}

type History struct {
	gorm.Model
	ID 			   	   uint   `gorm:"primaryKey"`
	UserID			   uint   `json:"user_id" gorm:"not null"`
	BookID			   uint   `json:"book_id" gorm:"not null"`
}

type Cart struct {
	gorm.Model
	ID 			   	   uint   `gorm:"primaryKey"`
	UserID			   uint   `json:"user_id" gorm:"not null"`
	BookID			   uint   `json:"book_id" gorm:"not null"`
}

type Review struct {
	gorm.Model
	ID 			   	   uint   `gorm:"primaryKey"`
	UserID			   uint   `json:"user_id" gorm:"not null"`
	BookID			   uint   `json:"book_id" gorm:"not null"`
	Review			   string `json:"review"`
	Rating			   uint   `json:"rating"`
}