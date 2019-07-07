package middleware

import (
	"net/http"
	"strings"

	"lenslocked.com/context"

	"lenslocked.com/models"
)

//RequireUser will make sure the user has logged in
type RequireUser struct {
	User
}

//Apply shall execute the middleware wherever it is called
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

//User will set the user context accross the app!!!
type User struct {
	models.UserService
}

//Apply shall execute the middleware wherever it is called
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

//ApplyFn is the same as apply, but uses HandlerFunc instead of Handler
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		//do not read the db for a static asset...
		if strings.HasPrefix(path, "/assets/") ||
			strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}

//ApplyFn is the same as apply, but uses HandlerFunc instead of Handler
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	ourHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		next(w, r)
	})
	return mw.User.Apply(ourHandler)
}
