package main

import 	(
	"net/http"
	 "fmt"
 )

func handlerFunc(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "text/html")
	if(r.URL.Path == "/"){
		fmt.Fprint(w, "<h1>Hey dog lover! Welcome!</h1>")
	} else if(r.URL.Path == "/contact") {
		fmt.Fprint(w, "To get in touch, please drop us a line at: <a href =\"mailto:support@lenslocked.com\"> support@lenslocked.com</a>")
	} else {
		fmt.Fprint(w, "Ooops... it appears we don't have this page arround.")
	}
}

func main()  {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":8080", nil)
}
