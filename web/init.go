package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/pkger"
	"gopkg.in/boj/redistore.v1"
	"html/template"
	"io/ioutil"
	"mega_bot/config"
	"net/http"
	"os"
	"strings"
)

type contextKey int

type templateAlert struct {
	Header string
	Text   string
}

type templateCommon struct {
	AlertError   templateAlert
	AlertSuccess templateAlert
	AlertWarn    templateAlert

	PageTitle string
}

const SessionKey contextKey = 0
const GothUserKey contextKey = 1

var (
	logger    *loggo.Logger
	store     *redistore.RediStore
	templates *template.Template
)

func Close() {
	store.Close()
}

func Init(conf *config.Config) error {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	// Load Templates
	templateDir := pkger.Include("/web/templates")
	t, err := compileTemplates(templateDir)
	if err != nil {
		return err
	}
	templates = t

	// Fetch new store.
	store, err = redistore.NewRediStoreWithDB(10, "tcp", conf.RedisAddress, conf.RedisPassword, "1", []byte(conf.CookieSecret))
	if err != nil {
		return err
	}

	// Init goth
	gothic.Store = store
	goth.UseProviders(
		discord.New(
			conf.DiscordKey,
			conf.DiscordSecret,
			fmt.Sprintf("https://%s/auth/discord/callback", conf.ExtHostname),
			discord.ScopeIdentify, discord.ScopeEmail),
	)

	// Setup Router
	r := mux.NewRouter()

	// Unprotected Pages
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(pkger.Dir("/web/static"))))
	r.HandleFunc("/auth/{provider}", GetAuthProvider).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", GetAuthProviderCallback).Methods("GET")
	r.HandleFunc("/login", GetLogin).Methods("GET")

	// Protected Pages
	r.HandleFunc("/", GetHome).Methods("GET")

	go func() {
		err := http.ListenAndServe(":9000", r)
		if err != nil {
			logger.Errorf("Could not start web server %s", err.Error())
		}
	}()

	return nil
}

func compileTemplates(dir string) (*template.Template, error) {
	const fun = "compileTemplates"
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
