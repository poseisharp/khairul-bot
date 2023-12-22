package services

import "github.com/bwmarrin/discordgo"

type DiscordService struct {
	s *discordgo.Session
}

func NewDiscordService(s *discordgo.Session) *DiscordService {
	return &DiscordService{s: s}
}

func (d *DiscordService) SendTextMessage(channelID string, message string) error {
	_, err := d.s.ChannelMessageSend(channelID, message)
	return err
}
