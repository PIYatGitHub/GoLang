package middleware

import (
	"fmt"
	"net/http"

	"../context"

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
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		fmt.Println("User Found: ... ", user)
		next(w, r) // call next in line -- done!

	})
}
