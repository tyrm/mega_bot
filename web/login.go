package web

import "net/http"

type LoginTemplate struct {
	templateCommon
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	tmplVars := LoginTemplate{}
	tmplVars.PageTitle = "Login"

	err := templates.ExecuteTemplate(w, "login", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
