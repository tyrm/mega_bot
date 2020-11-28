package web

import (
	"encoding/gob"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/pkger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/boj/redistore.v1"
	"html/template"
	"mega_bot/config"
	"mega_bot/models"
	"net/http"
	"regexp"
	"time"
)

const SessionKey contextKey = 0
const UserKey contextKey = 1
const LocalizerKey contextKey = 2

var (
	langBundle *i18n.Bundle
	logger     *loggo.Logger
	store      *redistore.RediStore
	templates  *template.Template
)

type contextKey int

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

	// Load Languages
	bundle, err := compileLanguages()
	if err != nil {
		return err
	}
	langBundle = bundle

	// Fetch new store.
	store, err = redistore.NewRediStoreWithDB(10, "tcp", conf.RedisAddress, conf.RedisPassword, "1", []byte(conf.CookieSecret))
	if err != nil {
		return err
	}

	// Register models for GOB
	gob.Register(models.User{})

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
	r.Use(Middleware)

	// Error Pages
	r.NotFoundHandler = NotFoundHandler()
	r.MethodNotAllowedHandler = MethodNotAllowedHandler()

	// Static Files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(pkger.Dir("/web/static"))))

	// Unprotected Pages
	r.HandleFunc("/auth/{provider}", GetAuthProvider).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", GetAuthProviderCallback).Methods("GET")
	r.HandleFunc("/login", GetLogin).Methods("GET")
	r.HandleFunc("/logout", GetLogout).Methods("GET")

	// Protected Pages
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(MiddlewareRequireAuth)
	protected.HandleFunc("/", GetHome).Methods("GET")
	protected.HandleFunc("/purgatory", GetPurgatory).Methods("GET")
	protected.HandleFunc("/responder", GetResponder).Methods("GET")
	protected.HandleFunc("/responder/add", GetResponderAdd).Methods("GET")
	protected.HandleFunc("/responder/edit/{responder}", GetResponderEdit).Methods("GET")

	go func() {
		srv := &http.Server{
			Handler:      r,
			Addr:         ":9000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			logger.Errorf("Could not start web server %s", err.Error())
		}
	}()

	return nil
}

// privates
func compileLanguages() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Files muse be listed with Include for pkger to pull them in
	files := map[string]string{
		"active.es.toml": pkger.Include("/active.es.toml"),
	}

	for filename, file := range files {
		langFile, err := pkger.Open(file)
		if err != nil {
			return nil, err
		}
		defer langFile.Close()

		fileinfo, err := langFile.Stat()
		if err != nil {
			return nil, err
		}

		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)

		_, err = langFile.Read(buffer)
		if err != nil {
			return nil, err
		}

		bundle.MustParseMessageFileBytes(buffer, filename)
	}

	return bundle, nil
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func roundUp(f float64) uint64 {
	fint := uint64(f)

	if f > float64(fint) {
		fint++
	}

	logger.Debugf("roundUp: %f int: %d", f, fint)
	return fint
}
