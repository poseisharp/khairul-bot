package features

import "github.com/bwmarrin/discordgo"

type FeatureCommand interface {
	DiscordCommand() *discordgo.ApplicationCommand
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error
}
