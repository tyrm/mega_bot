package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"mega_bot/models"
	"net/http"
	"reflect"
)

type LoginTemplate struct {
	templateCommon
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// Init template variables
	tmplVars := &LoginTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HomeTitle",
			Description: "Title of the home page.",
			Other: "Home",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)
		logger.Debugf("user(%s) %#v", reflect.TypeOf(user), user)
	}


	err = templates.ExecuteTemplate(w, "login", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
