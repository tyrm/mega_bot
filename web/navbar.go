package web

import "regexp"
import "github.com/jinzhu/copier"

var NavbarTemplate = []templateNavbarNode{
	{
		Text:     "Home",
		MatchStr: "^/$",
		FAIcon:   "home",
		URL:      "/",
	},
	{
		Text:     "Responder",
		MatchStr: "^/responder/.*$",
		FAIcon:   "comment-alt",
		URL:      "/responder/",
	},
	{
		Text:   "Admin",
		FAIcon: "hammer",
		URL:    "#",
		Children: []*templateNavbarNode{
			{
				Text:     "Job Runner",
				MatchStr: "^/web/admin/jobrunner/.*$",
				FAIcon:   "clock",
				URL:      "/web/admin/jobrunner/",
			},
			{
				Text:     "Oauth Clients",
				MatchStr: "^/web/admin/oauth-clients/.*$",
				FAIcon:   "desktop",
				URL:      "/web/admin/oauth-clients/",
				Disabled: true,
			},
			{
				Text:     "Registry",
				MatchStr: "^/web/admin/registry/.*$",
				FAIcon:   "book",
				URL:      "/web/admin/registry/",
			},
			{
				Text:     "Users",
				MatchStr: "^/web/admin/users/.*$",
				FAIcon:   "user",
				URL:      "/web/admin/users/",
			},
			{
				Text:     "Something else here",
				FAIcon:   "paw",
				URL:      "#",
				Disabled: true,
			},
		},
	},
}

var NavbarAdminTemplate = templateNavbarNode{
	Text:   "Admin",
	FAIcon: "hammer",
	URL:    "#",
	Children: []*templateNavbarNode{
		{
			Text:     "Job Runner",
			MatchStr: "^/web/admin/jobrunner/.*$",
			FAIcon:   "clock",
			URL:      "/web/admin/jobrunner/",
		},
		{
			Text:     "Oauth Clients",
			MatchStr: "^/web/admin/oauth-clients/.*$",
			FAIcon:   "desktop",
			URL:      "/web/admin/oauth-clients/",
			Disabled: true,
		},
		{
			Text:     "Registry",
			MatchStr: "^/web/admin/registry/.*$",
			FAIcon:   "book",
			URL:      "/web/admin/registry/",
		},
		{
			Text:     "Users",
			MatchStr: "^/web/admin/users/.*$",
			FAIcon:   "user",
			URL:      "/web/admin/users/",
		},
		{
			Text:     "Something else here",
			FAIcon:   "paw",
			URL:      "#",
			Disabled: true,
		},
	},
}

func makeNavbar(path string) (navbar *[]templateNavbarNode) {
	var newNavbar []templateNavbarNode
	err := copier.Copy(&newNavbar, &NavbarTemplate)
	if err != nil {
		logger.Errorf("could not copy navbar template: %s", err.Error())
	}


	for i := 0; i < len(newNavbar); i++ {
		if newNavbar[i].MatchStr != "" {
			match, err := regexp.MatchString(newNavbar[i].MatchStr, path)
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
					subMatch, err := regexp.MatchString(newNavbar[i].Children[j].MatchStr, path)
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