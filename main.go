package main

import (
	"fmt"
	"net/http"

	"./models"

	"../lenslocked.com/controllers"
	"github.com/gorilla/mux"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	services, err := models.NewServices(psqlInfo)
	must(err)
	// defer us.Close()
	// us.DestructiveReset()
	// -- but it works bad with fresh, so we use AutoMigrate
	// us.AutoMigrate()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUser(services.User)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
