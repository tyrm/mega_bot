package web

import (
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
	AlertError   templateAlert
	AlertSuccess templateAlert
	AlertWarn    templateAlert

	NavBar       *[]templateNavbarNode
	PageTitle    string
	User         *models.User
}

func (t *templateCommon) SetNavbar(n *[]templateNavbarNode) {
	t.NavBar = n
	return
}

func (t *templateCommon) SetUser(u *models.User) {
	t.User = u
	return
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
	SetNavbar(n *[]templateNavbarNode)
	SetUser(u *models.User)
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

func initTemplate(r *http.Request, tmpl templateVars) error {
	//us := r.Context().Value(SessionKey).(*sessions.Session)

	// add navbar
	tmpl.SetNavbar(makeNavbar(r.URL.Path))

	// add user to template
	if r.Context().Value(UserKey) != nil {
		user := r.Context().Value(UserKey).(*models.User)
		tmpl.SetUser(user)
	}

	// add alert error to template



	return nil
}