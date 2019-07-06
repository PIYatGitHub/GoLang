package models

import (
	"github.com/jinzhu/gorm"
	// added it not to get confused as of what is needed to run this...
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//NewServices will init all services with a single DB connection
func NewServices(dialect, connectionInfo string) (*Services, error) {
	// TODO: config this
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
}

// Services defines all we have -- for start it is Gallery and User services
type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

// Close will terminate the connection to the DB!
func (s *Services) Close() error {
	return s.db.Close()
}

//DestructiveReset deletes all available tables. NEVER EVER RUN IN PROD!!!!!
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

//AutoMigrate will migrat the db tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}
