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

//NewGalleryService will init the gallery service and will make it available
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

type galleryValidator struct {
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

type galleryService struct {
	GalleryDB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
