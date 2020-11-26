package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
	"net/http"
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
	tmplVars.PageTitle, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginTitle",
			Description: "Title of the login page.",
			Other:       "Login",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.Header, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginHeader",
			Description: "Header of the login page.",
			Other:       "Please sign in",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.LabelEmailAddress, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginLabelEmailAddress",
			Description: "Email address field label in the login form.",
			Other:       "Email address",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.LabelPassword, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginLabelPassword",
			Description: "Password field label in the login form.",
			Other:       "Password",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.ButtonSignIn, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginButtonSignIn",
			Description: "Sign In button in the login form.",
			Other:       "Sign In",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.ButtonForgotPassword, err = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginButtonForgotPassword",
			Description: "Forgot Password button in the login form.",
			Other:       "Forgot Password",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	buttonLoginDiscord, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "LoginButtonLoginDiscord",
			Description: "Login with Discord button in the login form.",
			Other:       "Login with {{.Icon}} Discord",
		},
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
