package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"mega_bot/models"
)

func handleGuildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {
	logger.Tracef("GuildCreate[%#v]", m.Guild)
}

func handleGuildDelete(s *discordgo.Session, m *discordgo.GuildDelete) {
	logger.Tracef("GuildDelete[%#v]", m.Guild)
}

func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	logger.Tracef("Message[ID: %s ChannelID: %s GuildID: %s Type: %v Content: %s]", m.ID, m.ChannelID, m.GuildID, m.Type, m.Content)

	isDM := false
	if m.GuildID == "" {
		isDM = true
	}

	rreq := models.ResponderRequest{
		DirectMessage: isDM,
		MeString: fmt.Sprintf("<@!%s>", s.State.User.ID),
		Message: m.Content,
		Service: "discord",
		ResponseTarget: m.ChannelID,
	}

	go func() {
		*chanResponderRequest <- &rreq
	}()

	//r := []rune(m.Content)
	//l1ogger.Debugf("%d", r[0])
}

func handlePresenceUpdate(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	logger.Tracef("PresenceUpdate[%#v]", m)
}

func handlePresencesReplace(s *discordgo.Session, m *discordgo.PresencesReplace) {
	logger.Tracef("PresencesReplace[%#v]", m)
}

func handleRateLimit(s *discordgo.Session, m *discordgo.RateLimit) {
	logger.Warningf("RateLimit[%#v]", m)
}

func handleRelationshipAdd(s *discordgo.Session, m *discordgo.RelationshipAdd) {
	logger.Tracef("RelationshipAdd[%#v]", m)
}

func handleRelationshipRemove(s *discordgo.Session, m *discordgo.RelationshipRemove) {
	logger.Tracef("RelationshipRemove[%#v]", m)
}