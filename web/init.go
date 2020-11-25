package web

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/pkger"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"gopkg.in/boj/redistore.v1"
	"html/template"
	"mega_bot/config"
	"mega_bot/models"
	"net/http"
	"time"
)

const SessionKey contextKey = 0
const UserKey contextKey = 1

var (
	logger    *loggo.Logger
	minifier  *minify.M
	store     *redistore.RediStore
	templates *template.Template
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

	// Load Minify
	minifier = minify.New()
	minifier.AddFunc("text/html", html.Minify)

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
