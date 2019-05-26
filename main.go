package main

import (
	"net/http"

	"../lenslocked.com/controllers"
	"github.com/gorilla/mux"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUser()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":8080", r)
}
