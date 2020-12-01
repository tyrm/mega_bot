package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"mega_bot/models"
	"net/http"
	"strconv"
	"strings"
)

type ResponderTemplate struct {
	templateCommon

	RMs          *[]models.ResponderMatcher
	RMPagination *[]templatePaginationItems

	// i18n
	ButtonAdd string
	Header    string
}

type ResponderFormTemplate struct {
	templateCommon
	Breadcrumbs *[]templateBreadcrumb

	RM *models.ResponderMatcher

	// i18n
	ButtonSubmit       string
	Header             string
	LabelAlwaysRespond string
	LabelDescription   string
	LabelEnabled       string
	LabelID            string
	LabelMatchRegex    string
	LabelResponse      string
}

func GetResponder(w http.ResponseWriter, r *http.Request) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// init template variables
	tmplVars := &ResponderTemplate{}
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

	// get display count
	displayCount := 20
	hrefCount := 0
	if qCount, ok := r.URL.Query()["count"]; ok {
		if len(qCount[0]) >= 1 {
			uCount, err := strconv.ParseUint(qCount[0], 10, 32)
			if err != nil {
				logger.Debugf("invalid count: %s", qCount[0])
			} else {
				displayCount = int(uCount)
				hrefCount = int(uCount)
			}
		}
	}

	pages := roundUp(float64(rmCount) / float64(displayCount))
	logger.Debugf("pages: %d", pages)

	// get display page
	displayPage := 1
	if qPage, ok := r.URL.Query()["page"]; ok {
		if len(qPage[0]) >= 1 {
			uPage, err := strconv.ParseUint(qPage[0], 10, 32)
			if err != nil {
				logger.Debugf("invalid page: %s", qPage[0])
			} else {
				displayPage = int(uPage)
			}
		}
	}

	// mage pagination
	tmplVars.RMPagination = makePagination(displayPage, pages, "/responder", hrefCount)

	// get list of responders
	rms, err := models.ReadResponderMatchersPage(displayPage-1, displayCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplVars.RMs = rms

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

	err = templates.ExecuteTemplate(w, "responder", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}

}

func GetResponderAdd(w http.ResponseWriter, r *http.Request) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// i18n
	locAddResponder, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textAddResponder, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	returnResponderForm(w, r, nil, locAddResponder)

}

func GetResponderEdit(w http.ResponseWriter, r *http.Request) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// get responder
	vars := mux.Vars(r)
	if !isValidUUID4(vars["responder"]) {
		returnErrorPage(w, r, http.StatusBadRequest, "invalid id format")
		return
	}
	rm, err := models.ReadResponderMatcher(vars["responder"])
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if rm == nil {
		returnErrorPage(w, r, http.StatusNotFound, fmt.Sprintf("responder not found: %s", vars["responder"]))
		return
	}

	// i18n
	locEditResponder, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textEditResponder, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	returnResponderForm(w, r, rm, locEditResponder)
}

func PostResponderAdd(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseForm()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	rm := models.ResponderMatcher{
		AlwaysRespond: false,
		Description:   r.Form.Get("description"),
		Enabled:       false,
		MatcherString: r.Form.Get("regex"),
		Response:      r.Form.Get("response"),
	}

	if r.Form.Get("enabled") == "on" {
		rm.Enabled = true
	}

	if r.Form.Get("always-respond") == "on" {
		rm.AlwaysRespond = true
	}

	// validate regex
	err = rm.BuildRE()
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = models.CreateResponderMatcher(&rm)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Responder added"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, "/responder", http.StatusFound)

}

func PostResponderEdit(w http.ResponseWriter, r *http.Request) {
	// get responder
	vars := mux.Vars(r)
	if !isValidUUID4(vars["responder"]) {
		returnErrorPage(w, r, http.StatusBadRequest, "invalid id format")
		return
	}

	// parse form data
	err := r.ParseForm()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	rm := models.ResponderMatcher{
		ID:            vars["responder"],
		AlwaysRespond: false,
		Description:   r.Form.Get("description"),
		Enabled:       false,
		MatcherString: r.Form.Get("regex"),
		Response:      r.Form.Get("response"),
	}

	if r.Form.Get("enabled") == "on" {
		rm.Enabled = true
	}

	if r.Form.Get("always-respond") == "on" {
		rm.AlwaysRespond = true
	}

	// validate regex
	err = rm.BuildRE()
	if err != nil {
		returnErrorPage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = models.UpdateResponderMatcher(&rm)
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	us := r.Context().Value(SessionKey).(*sessions.Session)
	us.Values["page-alert-success"] = templateAlert{Text: "Responder updated"}
	err = us.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect home
	http.Redirect(w, r, "/responder", http.StatusFound)
}

func returnResponderForm(w http.ResponseWriter, r *http.Request, rm *models.ResponderMatcher, actionText string) {
	// get localizer
	localizer := r.Context().Value(LocalizerKey).(*i18n.Localizer)

	// Init template variables
	tmplVars := &ResponderFormTemplate{}
	err := initTemplate(w, r, tmplVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set
	tmplVars.RM = rm

	// i18n
	tmplVars.PageTitle = strings.Title(actionText)
	tmplVars.Header = strings.Title(actionText)
	tmplVars.ButtonSubmit = strings.Title(actionText)

	tmplVars.LabelAlwaysRespond, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textAlwaysRespond})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelAlwaysRespond = strings.Title(tmplVars.LabelAlwaysRespond)

	tmplVars.LabelDescription, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textDescription, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelDescription = strings.Title(tmplVars.LabelDescription)

	tmplVars.LabelEnabled, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textEnabled})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelEnabled = strings.Title(tmplVars.LabelEnabled)

	tmplVars.LabelID, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textID, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	tmplVars.LabelMatchRegex, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textMatchRegex, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelMatchRegex = strings.Title(tmplVars.LabelMatchRegex)

	tmplVars.LabelResponse, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textResponse, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}
	tmplVars.LabelResponse = strings.Title(tmplVars.LabelResponse)

	locResponder, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &textResponder, PluralCount: 1})
	if err != nil {
		logger.Warningf("missing translation: %s", err.Error())
	}

	// breadcrumbs
	breadcrumbs := []templateBreadcrumb{
		{
			HRef: "/responder",
			Text: strings.Title(locResponder),
		},
		{
			Text: strings.Title(actionText),
		},
	}
	tmplVars.Breadcrumbs = &breadcrumbs

	err = templates.ExecuteTemplate(w, "responder_form", tmplVars)
	if err != nil {
		logger.Errorf("could not render home template: %s", err.Error())
	}
}
