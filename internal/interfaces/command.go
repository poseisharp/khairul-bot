package interfaces

import "github.com/bwmarrin/discordgo"

type FeatureCommand interface {
	DiscordCommand() *discordgo.ApplicationCommand
	HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error
}
