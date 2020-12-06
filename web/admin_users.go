package web

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"mega_bot/models"
	"net/http"
	"strings"
)

type UsersTemplate struct {
	templateCommon

	Users      *[]models.UserListItem
	Pagination *[]templatePaginationItems

	// i18n
	ButtonAdd string
	Header    string
}

func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// init template variables
	tmplVars := &UsersTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get count of RMs
	rmCount, err := models.CountResponderMatchers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pr := getPageFromURL(r, rmCount, 20)

	// mage pagination
	tmplVars.Pagination = makePagination(pr.CurrentPage, pr.Pages, "/responder", pr.HrefCount)

	// get list of responders
	userlist, err := models.ReadUserListPage(pr.CurrentPage-1, pr.DisplayCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplVars.Users = userlist

	// i18n
	locResponder, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textResponder, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.PageTitle = strings.Title(locResponder)
	tmplVars.Header = strings.Title(locResponder)

	tmplVars.ButtonAdd, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textAdd})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.ButtonAdd = strings.Title(tmplVars.ButtonAdd)

	err = templates.ExecuteTemplate(w, "admin_users", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}
