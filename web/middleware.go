package web

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

		// Init Localizer
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(langBundle, lang, accept)
		ctx = context.WithValue(ctx, LocalizerKey, localizer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MiddlewareRequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Context().Value(UserKey) == nil {
			us := r.Context().Value(SessionKey).(*sessions.Session)

			// Save current page
			us.Values["login-redirect"] = r.URL.Path
			err := us.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// redirect to login
			returnErrorPage(w, r, http.StatusUnauthorized, "")
			return
		} else {
			user := r.Context().Value(UserKey).(*models.User)

			if !user.Authorized && r.URL.Path != "/purgatory" {
				http.Redirect(w, r, "/purgatory", http.StatusFound)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
