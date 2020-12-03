//go:generate pkger
package main

import (
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"gopkg.in/alexcesaro/statsd.v2"
	"mega_bot/config"
	"mega_bot/discord"
	"mega_bot/models"
	"mega_bot/responder"
	"mega_bot/web"
	"os"
	"os/signal"
	"syscall"
)

var logger *loggo.Logger

func main() {
	conf := config.CollectConfig()

	// Init Logging
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	err := loggo.ConfigureLoggers(conf.LoggerConfig)
	if err != nil {
		logger.Errorf("Error configuring Logger: %s", err.Error())
		return
	}
	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		logger.Errorf("Error configuring Color Logger: %s", err.Error())
		return
	}

	logger.Infof("Starting Mega Bot")

	// init statsd
	sdc, err := statsd.New(
		statsd.Address(conf.StatsDAddress),
		statsd.Prefix("megabot"),
		statsd.TagsFormat(statsd.InfluxDB),
	) // Connect to the UDP port 8125 by default.

	if err != nil {
		logger.Warningf("statsd: %s", err.Error())
	}
	defer sdc.Close()
	StartStatsSender(sdc)

	// Communication Channels
	var chanResponderRequest chan *models.ResponderRequest
	chanResponderRequest = make(chan *models.ResponderRequest)

	// Init internals
	err = models.Init(conf, sdc)
	if err != nil {
		logger.Errorf("unable to connect to database: %s", err.Error())
		return
	}

	err = responder.Init(conf, &chanResponderRequest)
	if err != nil {
		logger.Errorf("unable to connect to database: %s", err.Error())
		return
	}

	err = web.Init(conf)
	if err != nil {
		logger.Errorf("unable to start webserver: %s", err.Error())
		return
	}

	// Init Bot Connections
	if conf.DiscordToken != "" {
		err = discord.Init(conf, &chanResponderRequest)
		if err != nil {
			logger.Errorf("unable to connect to discord: %s", err.Error())
		}
	}

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	web.Close()
}
