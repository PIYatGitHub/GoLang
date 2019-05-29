package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService, error) {
	const (
		host   = "localhost"
		port   = 5432
		user   = "postgres"
		dbname = "lenslocked_test"
	)
	// ^^^^ DO NOT FORGET:: create the DB manually! ^^^^
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)
	// clear the users table before we test!
	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "Mika Hackinen",
		Email: "mika@hackinen.io",
	}
	err = us.Create(&user)

	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID to be a positive int, received a 0")
	}

	if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected to have created within 5 seconds...")
	}

	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected to have updated within 5 seconds...")
	}

}
func TestUserByID(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "Mika Hackinen",
		Email: "mika@hackinen.io",
	}
	err = us.Create(&user)

	if err != nil {
		t.Fatal(err)
	}

	_, err = us.ByID(1)
	if err != nil {
		t.Errorf("Expected to find the user at position 1, but failed...")
	}
}
