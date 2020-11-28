package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
	"net/http"
	"strings"
)

type LoginTemplate struct {
	templateCommon

	Header               string
	LabelEmailAddress    string
	LabelPassword        string
	ButtonSignIn         string
	ButtonForgotPassword string
	ButtonLoginDiscord   template.HTML
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

	// change CSS sheet
	tmplVars.HeadCSS = &[]templateHeadLink{
		{
			HRef: "/static/css/login.css",
			Rel: "stylesheet",
		},
	}

	// disable navbar
	tmplVars.NavBarEnabled = false

	// custom body css
	tmplVars.BodyClass = "text-center"

	// i18n
	tmplVars.PageTitle, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textLogin})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.PageTitle = strings.Title(tmplVars.PageTitle)

	tmplVars.Header, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textSignInAsk})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.LabelEmailAddress, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textEmailAddress})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelEmailAddress = strings.Title(tmplVars.LabelEmailAddress)

	tmplVars.LabelPassword, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textPassword})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelPassword = strings.Title(tmplVars.LabelPassword)

	tmplVars.ButtonSignIn, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textSignIn})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.ButtonSignIn = strings.Title(tmplVars.ButtonSignIn)

	tmplVars.ButtonForgotPassword, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textForgotPassword})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.ButtonForgotPassword = strings.Title(tmplVars.ButtonForgotPassword)

	buttonLoginDiscord, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &textLoginDiscord,
		TemplateData: map[string]interface{}{
			"Icon": "<i class=\"fab fa-discord\"></i>",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.ButtonLoginDiscord = template.HTML(buttonLoginDiscord) // preserve html

	err = templates.ExecuteTemplate(w, "login", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
