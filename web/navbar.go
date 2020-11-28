package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"regexp"
	"strings"
)

func makeNavbar(r *http.Request) (navbar *[]templateNavbarNode) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// i18n stuff
	natbarTextHome, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textWebHome})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	natbarTextHome = strings.Title(natbarTextHome)

	natbarTextResponder, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textResponder, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	natbarTextResponder = strings.Title(natbarTextResponder)

	natbarTextAdmin, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textAdmin, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	natbarTextAdmin = strings.Title(natbarTextAdmin)

	natbarTextUsers, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textUser})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	natbarTextUsers = strings.Title(natbarTextUsers)

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
			MatchStr: "^/responder.*$",
			FAIcon:   "comment-alt",
			URL:      "/responder",
		},
		{
			Text: natbarTextAdmin,
			FAIcon: "hammer",
			URL:    "#",
			Children: []*templateNavbarNode{
				{
					Text:     natbarTextUsers,
					MatchStr: "^/admin/users/.*$",
					FAIcon:   "user",
					URL:      "/admin/users/",
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
