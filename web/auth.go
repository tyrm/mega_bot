package web

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"mega_bot/models"
	"net/http"
)

func GetAuthProvider(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		authProviderCallback(w, r, gothUser)
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

	authProviderCallback(w, r, user)
	return
}

func authProviderCallback(w http.ResponseWriter, r *http.Request, gu goth.User) {
	us := r.Context().Value(SessionKey).(*sessions.Session)
	var loggedInUser *models.User
	if r.Context().Value(UserKey) != nil {
		loggedInUser = r.Context().Value(UserKey).(*models.User)
	}

	logger.Tracef("loggedInUser: %#v", loggedInUser)
	logger.Tracef("gothUser: %#v", gu)

	ca, err := models.ReadConnectedAccount(gu.Provider, gu.UserID)
	if err != nil {
		msg := fmt.Sprintf("could not read connected account: %s", err.Error())
		logger.Errorf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	logger.Tracef("ca: %#v", ca)

	if loggedInUser != nil && ca != nil && ca.UserID == loggedInUser.ID {
		// user is logged in, the account is already connected, everything is fine. go home
		us.Values["page-alert-success"] = &templateAlert{Text: "This account is already connected"}
		err := us.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if loggedInUser == nil && ca != nil {
		// user is not logged in, but their account is linked, login the user
		u, err := models.ReadUser(ca.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		us.Values["user"] = u
		err = us.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if loggedInUser != nil && ca == nil {
		// user is logged in, account not connected. connect account, go home
		newCa := models.ConnectedAccount{
			Provider: gu.Provider,
			ProviderUserID: gu.UserID,
			UserID: loggedInUser.ID,
		}

		err := models.CreateConnectedAccount(&newCa)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		us.Values["page-alert-success"] = &templateAlert{Text: "Account connected."}
		err = us.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect home
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if loggedInUser == nil && ca == nil {
		// user is not logged in and there is no connected account

		// check for existing user
		u, err := models.ReadUserByEmail(gu.Email)
		if err != nil {
			msg := fmt.Sprintf("could not read user: %s", err.Error())
			logger.Errorf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if u != nil {
			// there is an existing user with this email address we cannot create a new account. ask to link
			us.Values["page-alert-warn"] = &templateAlert{
				Header: "Oops",
				Text: "An account already exists for that email address. Are you triyng to link an account? Please " +
					"login first.",
			}
			err := us.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// create user account
		newUser := models.User{}
		newUser.Email = gu.Email
		if gu.NickName != "" {
			newUser.Nick.String = gu.NickName
			newUser.Nick.Valid = true
		} else if gu.Name != "" {
			newUser.Nick.String = gu.Name
			newUser.Nick.Valid = true
		}
		newUser.Authorized = false
		newUser.Admin = false

		err = models.CreateUser(&newUser)
		if err != nil {
			msg := fmt.Sprintf("could not create new user: %s", err.Error())
			logger.Errorf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		logger.Debugf("created new user: %#v", newUser)

		// login user
		us.Values["user"] = &newUser

		// connect account
		newCa := models.ConnectedAccount{
			Provider: gu.Provider,
			ProviderUserID: gu.UserID,
			UserID: newUser.ID,
		}

		err = models.CreateConnectedAccount(&newCa)
		if err != nil {
			msg := fmt.Sprintf("could not create new connected account: %s", err.Error())
			logger.Errorf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		err = us.Save(r, w)
		if err != nil {
			msg := fmt.Sprintf("could not save session: %s", err.Error())
			logger.Errorf(msg)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Error(w, "unknown state", http.StatusInternalServerError)
	return
}
