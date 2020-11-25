package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type HomeTemplate struct {
	templateCommon
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	// localizer
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(langBundle, lang, accept)

	tmplVars := &HomeTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmplVars.PageTitle = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TitleHome",
			Description: "Title of the home page.",
			Other: "Home",
		},
	})

	err = templates.ExecuteTemplate(w, "home", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
