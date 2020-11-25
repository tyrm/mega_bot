package web

import (
	"context"
	"github.com/gorilla/sessions"
	"mega_bot/models"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Init Session
		us, err := store.Get(r, "megabot")
		if err != nil {
			logger.Infof("got %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), SessionKey, us)

		// Retrieve our user and type-assert it
		val := us.Values["user"]
		var user = models.User{}
		var ok bool
		if user, ok = val.(models.User); ok {
			ctx = context.WithValue(ctx, UserKey, &user)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MiddlewareRequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		us := r.Context().Value(SessionKey).(*sessions.Session)
		if r.Context().Value(UserKey) == nil {
			// Save current page
			us.Values["login-redirect"] = r.URL.Path
			err := us.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// redirect to login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
