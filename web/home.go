package web

import "net/http"

func GetHome(w http.ResponseWriter, r *http.Request) {

	err := templates.ExecuteTemplate(w, "home", nil)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
