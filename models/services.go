package models

import "github.com/jinzhu/gorm"

//NewServices will init all services with a single DB connection
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &Services{
		User: NewUserService(db),
		db:   db,
	}, nil
}

// Services defines all we have -- for start it is Gallery and User services
type Services struct {
	Gallery GalleryService
	User    UserService
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
