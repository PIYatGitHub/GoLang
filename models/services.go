package models

import "github.com/jinzhu/gorm"

//NewServices will init all services with a single DB connection
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &Services{}, nil
}

// Services defines all we have -- for start it is Gallery and User services
type Services struct {
	Gallery GalleryService
	User    UserService
}
