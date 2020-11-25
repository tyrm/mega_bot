package web

import (
	"mega_bot/models"
	"net/http"
	"reflect"
)

type LoginTemplate struct {
	templateCommon
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	tmplVars := LoginTemplate{}
	tmplVars.PageTitle = "Login"

	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)
		logger.Debugf("user(%s) %#v", reflect.TypeOf(user), user)
	}


	err := templates.ExecuteTemplate(w, "login", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
