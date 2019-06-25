package models

import "github.com/jinzhu/gorm"

// Gallery will represent the core of our app - what the user is able to see
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

//GalleryService is the available API
type GalleryService interface {
	GalleryDB
}

// GalleryDB holds the CRUD for galleries
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}
