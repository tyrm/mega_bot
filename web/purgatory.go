package web

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"mega_bot/models"
	"net/http"
	"strings"
)

func GetPurgatory(w http.ResponseWriter, r *http.Request) {
	// check user state
	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)

		authorized, err := user.HasOneOfRoles([]string{"administrator", "operator", "authorized"})
		if err != nil {
			returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if authorized {
			// redirect to home, user is authorized
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// Init template variables
	tmplVars := &SinglePageTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// hdie navbar contents
	tmplVars.NavBar = &[]templateNavbarNode{}

	locPurgatory, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textPurgatory, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.PageTitle = strings.Title(locPurgatory)
	tmplVars.Header = fmt.Sprintf("%s \U0001f47b", strings.Title(locPurgatory))

	purgatoryText, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textPurgatoryPageText})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.Paragraphs = []string{purgatoryText}

	err = templates.ExecuteTemplate(w, "single_page", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}