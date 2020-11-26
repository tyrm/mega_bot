package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"regexp"
)

func makeNavbar(r *http.Request) (navbar *[]templateNavbarNode) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// i18n stuff
	natbarTextHome, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "NavebarHome",
			Description: "Home button on Navbar",
			Other:       "Home",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	natbarTextResponder, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "NavbarResponder",
			Description: "Responder button on Navbar",
			Other:       "Responder",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	natbarTextAdmin, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "NavbarAdmin",
			Description: "Admin button on Navbar",
			Other:       "Admin",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	natbarTextAdminUser, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "NavbarAdminUsers",
			Description: "Users button on Navbar in the Admin Menu",
			Other:       "Users",
		},
	})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	// create navbar
	newNavbar := []templateNavbarNode{
		{
			Text: natbarTextHome,
			MatchStr: "^/$",
			FAIcon:   "home",
			URL:      "/",
		},
		{
			Text: natbarTextResponder,
			MatchStr: "^/responder/.*$",
			FAIcon:   "comment-alt",
			URL:      "/responder/",
		},
		{
			Text: natbarTextAdmin,
			FAIcon: "hammer",
			URL:    "#",
			Children: []*templateNavbarNode{
				{
					Text: natbarTextAdminUser,
					MatchStr: "^/web/admin/users/.*$",
					FAIcon:   "user",
					URL:      "/web/admin/users/",
				},
			},
		},
	}

	for i := 0; i < len(newNavbar); i++ {
		if newNavbar[i].MatchStr != "" {
			match, err := regexp.MatchString(newNavbar[i].MatchStr, r.URL.Path)
			if err != nil {
				logger.Errorf("makeNavbar:Error matching regex: %v", err)
			}
			if match {
				newNavbar[i].Active = true
			}

		}

		if newNavbar[i].Children != nil {
			for j := 0; j < len(newNavbar[i].Children); j++ {

				if newNavbar[i].Children[j].MatchStr != "" {
					subMatch, err := regexp.MatchString(newNavbar[i].Children[j].MatchStr, r.URL.Path)
					if err != nil {
						logger.Errorf("makeNavbar:Error matching regex: %v", err)
					}

					if subMatch {
						newNavbar[i].Active = true
						newNavbar[i].Children[j].Active = true
					}

				}

			}
		}
	}

	return &newNavbar
}
