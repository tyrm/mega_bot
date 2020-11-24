package web

import (
	"context"
	"net/http"
)

func MiddlewareProtected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Init Session
		us, err := store.Get(r, "megabot")
		if err != nil {
			logger.Infof("got %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), SessionKey, &us)


		rWithSignature := r.WithContext(ctx)
		next.ServeHTTP(w, rWithSignature)
	})
}
