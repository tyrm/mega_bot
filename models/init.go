package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/loggo"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	"github.com/patrickmn/go-cache"
	"github.com/rubenv/sql-migrate"
	"gopkg.in/alexcesaro/statsd.v2"
	"mega_bot/config"
)

var (
	client         *sqlx.DB
	logger         *loggo.Logger
	stats          *statsd.Client
	cHasOneOfRoles *cache.Cache
)

func Init(conf *config.Config, sdc *statsd.Client) error {
	// Init Logging
	newLogger := loggo.GetLogger("models")
	logger = &newLogger

	var err error
	client, err = sqlx.Connect("postgres", conf.DBEngine)
	if err != nil {
		return err
	}

	// Statsd
	stats = sdc.Clone(statsd.Prefix("models"))

	// Do Migration
	logger.Debugf("Loading Migrations")
	migrations := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/models/migrations"),
	}

	logger.Debugf("Applying Migrations")
	n, err := migrate.Exec(client.DB, "postgres", migrations, migrate.Up)
	if n > 0 {
		logger.Infof("Applied %d migrations!\n", n)
	}
	if err != nil {
		logger.Criticalf("Coud not migrate database: %s", err)
		return err
	}
	return nil
}
