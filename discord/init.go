package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/juju/loggo"
	"mega_bot/config"
	"mega_bot/models"
)

var (
	chanResponderRequest *chan *models.ResponderRequest
	client *discordgo.Session
	logger *loggo.Logger
)

func Init(conf *config.Config, crr *chan *models.ResponderRequest) error {
	// Init Logging
	newLogger := loggo.GetLogger("discord")
	logger = &newLogger

	logger.Debugf("Connecting to Discord")

	// Save channel pointer
	chanResponderRequest = crr

	// Create a new Discord session using the provided bot token.
	var err error
	client, err = discordgo.New("Bot " + conf.DiscordToken)
	if err != nil {
		return err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	client.AddHandler(handleGuildCreate)
	client.AddHandler(handleGuildDelete)
	client.AddHandler(handleMessageCreate)
	client.AddHandler(handlePresenceUpdate)
	client.AddHandler(handlePresencesReplace)
	client.AddHandler(handleRateLimit)
	client.AddHandler(handleRelationshipAdd)
	client.AddHandler(handleRelationshipRemove)

	// In this example, we only care about receiving message events.
	client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	// Open a websocket connection to Discord and begin listening.
	err = client.Open()
	if err != nil {
		return err
	}

	return nil
}
