package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type HomeTemplate struct {
	templateCommon
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// Init template variables
	tmplVars := &HomeTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textWebHome})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	err = templates.ExecuteTemplate(w, "home", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
