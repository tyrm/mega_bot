package web

import "net/http"

type HomeTemplate struct {
	templateCommon
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	tmplVars := HomeTemplate{}
	tmplVars.PageTitle = "Home"

	err := templates.ExecuteTemplate(w, "home", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
