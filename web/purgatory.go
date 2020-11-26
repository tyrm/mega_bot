package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"mega_bot/models"
	"net/http"
)

func GetPurgatory(w http.ResponseWriter, r *http.Request) {
	// check user state
	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)

		if user.Authorized {
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

	tmplVars.NavBar = &[]templateNavbarNode{}

	tmplVars.PageTitle, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PurgatoryTitle",
			Description: "Title of the purgatory page.",
			Other: "Purgatory",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.Header, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PurgatoryHeader",
			Description: "Header of the purgatory page.",
			Other: "Purgatory {{.GhostEmoji}}",
		},
		TemplateData: map[string]interface{}{
			"GhostEmoji": "\U0001f47b",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	purgatoryText, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PurgatoryText",
			Description: "Text of the purgatory page.",
			Other: "If you're seeing this page you're not approved to use MegaBot. Please contact the administrator.",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.Paragraphs = []string{purgatoryText}

	err = templates.ExecuteTemplate(w, "single_page", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}