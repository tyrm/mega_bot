package web

import (
	"context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/alexcesaro/statsd.v2"
	"mega_bot/models"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := stats.Clone(statsd.Tags("path", r.URL.Path, "method", r.Method)).NewTiming()

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

		startTime.Send("response")
	})
}

func MiddlewareRequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Context().Value(UserKey) == nil {
			// redirect to login
			returnErrorPage(w, r, http.StatusUnauthorized, "")
			return
		} else {
			user := r.Context().Value(UserKey).(*models.User)

			authorized, err := user.HasOneOfRoles([]string{"administrator", "operator", "authorized"})
			if err != nil {
				returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
				return
			}

			if !authorized && r.URL.Path != "/purgatory" {
				http.Redirect(w, r, "/purgatory", http.StatusFound)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
