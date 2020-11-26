package web

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/copier"
	"github.com/markbates/pkger"
	"html/template"
	"io/ioutil"
	"mega_bot/models"
	"net/http"
	"os"
	"strings"
)

type templateAlert struct {
	Header string
	Text   string
}

type templateCommon struct {
	HeadCSS          *[]templateHeadLink
	HeadFavicons     *[]templateHeadLink
	HeadFrameworkCSS *[]templateHeadLink

	AlertError   *templateAlert
	AlertSuccess *templateAlert
	AlertWarn    *templateAlert

	NavBar    *[]templateNavbarNode
	PageTitle string
	User      *models.User
}

func (t *templateCommon) SetAlertError(a *templateAlert) {
	t.AlertError = a
	return
}

func (t *templateCommon) SetAlertSuccess(a *templateAlert) {
	t.AlertSuccess = a
	return
}

func (t *templateCommon) SetAlertWarn(a *templateAlert) {
	t.AlertWarn = a
	return
}

func (t *templateCommon) SetHeadCSS(l *[]templateHeadLink) {
	t.HeadCSS = l
	return
}

func (t *templateCommon) SetHeadFavicons(l *[]templateHeadLink) {
	t.HeadFavicons = l
	return
}

func (t *templateCommon) SetHeadFrameworkCSS(l *[]templateHeadLink) {
	t.HeadFrameworkCSS = l
	return
}

func (t *templateCommon) SetNavbar(n *[]templateNavbarNode) {
	t.NavBar = n
	return
}

func (t *templateCommon) SetUser(u *models.User) {
	t.User = u
	return
}

type templateHeadLink struct {
	HRef        string
	Rel         string
	Integrity   string
	CrossOrigin string
	Sizes       string
	Type        string
}

type templateNavbarNode struct {
	Text     string
	URL      string
	MatchStr string
	FAIcon   string

	Active   bool
	Disabled bool

	Children []*templateNavbarNode
}

type templateVars interface {
	SetAlertError(a *templateAlert)
	SetAlertSuccess(a *templateAlert)
	SetAlertWarn(a *templateAlert)
	SetHeadCSS(l *[]templateHeadLink)
	SetHeadFavicons(l *[]templateHeadLink)
	SetHeadFrameworkCSS(l *[]templateHeadLink)
	SetNavbar(n *[]templateNavbarNode)
	SetUser(u *models.User)
}


type SinglePageTemplate struct {
	templateCommon
	Header string
	Paragraphs []string
}

func compileTemplates(dir string) (*template.Template, error) {
	tpl := template.New("")

	err := pkger.Walk(dir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".gohtml") {
			return nil
		}
		f, err := pkger.Open(path)
		if err != nil {
			logger.Errorf("could not open pkger path %s: %s", path, err.Error())
			return err
		}
		// Now read it.
		sl, err := ioutil.ReadAll(f)
		if err != nil {
			logger.Errorf("could not read pkger file %s: %s", path, err.Error())
		}

		// It can now be parsed as a string.
		_, err = tpl.Parse(string(sl))
		if err != nil {
			logger.Errorf("could not open parse template %s: %s", path, err.Error())
			return err
		}

		return nil
	})
	return tpl, err
}

func initTemplate(w http.ResponseWriter, r *http.Request, tmpl templateVars) error {
	us := r.Context().Value(SessionKey).(*sessions.Session)
	saveSession := false

	// add navbar
	tmpl.SetNavbar(makeNavbar(r))

	// add user
	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)
		tmpl.SetUser(user)
	}

	// add css
	var headFrameworkCSS []templateHeadLink
	err := copier.Copy(&headFrameworkCSS, &HeadFrameworkCSSTemplate)
	if err != nil {
		return err
	}
	tmpl.SetHeadFrameworkCSS(&headFrameworkCSS)

	var headFavicons []templateHeadLink
	err = copier.Copy(&headFavicons, &HeadFaviconsTemplate)
	if err != nil {
		return err
	}
	tmpl.SetHeadFavicons(&headFavicons)

	var headCSS []templateHeadLink
	err = copier.Copy(&headCSS, &HeadCSSTemplate)
	if err != nil {
		return err
	}
	tmpl.SetHeadCSS(&headCSS)

	// add alerts
	if us.Values["page-alert-error"] != nil {
		alert := us.Values["page-alert-error"].(*templateAlert)
		tmpl.SetAlertError(alert)

		us.Values["page-alert-error"] = nil
		saveSession = true

	}

	if us.Values["page-alert-success"] != nil {
		alert := us.Values["page-alert-success"].(*templateAlert)
		tmpl.SetAlertSuccess(alert)

		us.Values["page-alert-success"] = nil
		saveSession = true
	}

	if us.Values["page-alert-warn"] != nil {
		alert := us.Values["page-alert-warn"].(*templateAlert)
		tmpl.SetAlertWarn(alert)

		us.Values["page-alert-warn"] = nil
		saveSession = true
	}

	if saveSession {
		err := us.Save(r, w)
		if err != nil {
			logger.Warningf("initTemplate could not save session: %s", err.Error())
			return err
		}
	}

	return nil
}
