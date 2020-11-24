package web

import (
	"github.com/markbates/goth/gothic"
	"net/http"
)

func GetAuthProvider(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		logger.Tracef("gothUser: %#v", gothUser)
		// redirect home page if no login-redirect
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func GetAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Tracef("user: %#v", user)

	// redirect home page if no login-redirect
	http.Redirect(w, r, "/", http.StatusFound)
	return
}