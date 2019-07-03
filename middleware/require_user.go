package middleware

import (
	"fmt"
	"net/http"
	"time"

	"../models"
)

//RequireUser will make sure the user has logged in
type RequireUser struct {
	models.UserService
}

//Apply shall execute the middleware wherever it is called
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

//ApplyFn is the same as apply, but uses HandlerFunc instead of Handler
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if the user is logged in...
		t := time.Now()
		fmt.Println("Fake req started @:", t)
		next(w, r) // call next in line -- done!
		fmt.Println("Fake req ending at:", time.Since(t))
	})
}
