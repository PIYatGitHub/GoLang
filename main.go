package main

import 	(
	"net/http"
	 "fmt"
 )

func handlerFunc(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "text/html")
	if(r.URL.Path == "/"){
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "<h1>Hey dog lover! Welcome!</h1>")
	} else if(r.URL.Path == "/contact") {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "To get in touch, please drop us a line at: <a href =\"mailto:support@lenslocked.com\"> support@lenslocked.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Ooops... it appears we don't have this page arround. Imagine a 404 here... ")
	}
}

func main()  {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":8080", nil)
}
