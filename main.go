package main

import 	(
	"net/http"
	 "fmt"
 )

func handlerFunc(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main()  {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":8080", nil)
}
