package models

import (
	"github.com/jinzhu/gorm"
	// added it not to get confused as of what is needed to run this...
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ServicesConfig is the type of function you use to meke the
//config magic happen!
type ServicesConfig func(*Services) error

//WithGorm will con the DB connection.
func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

//WithLogMode will config the db log mode flag
func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

//WithUser will config the user service
func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

//WithGallery will config the gallery service
func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

//WithImage will config the image service
func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

//NewServices will init all services with a single DB connection
func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
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
