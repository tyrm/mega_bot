package discord

func SendMessage(msg, target string) error {
	_, err := client.ChannelMessageSend(target, msg)
	return err
}
