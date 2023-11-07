package set_timezone_feature

import (
	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/interfaces/features"
)

type SetTimezoneCommand struct {
	features.FeatureCommand

	discordCommand *discordgo.ApplicationCommand
}

var _ features.FeatureCommand = &SetTimezoneCommand{}

func New() *SetTimezoneCommand {
	return &SetTimezoneCommand{
		discordCommand: &discordgo.ApplicationCommand{
			Name:        "set_timezone",
			Description: "Set your timezone",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timezone",
					Description: "Your timezone",
					Required:    true,
				},
			},
		},
	}
}

func (p *SetTimezoneCommand) DiscordCommand() *discordgo.ApplicationCommand {
	return p.discordCommand
}

func (*SetTimezoneCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Code here
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Success",
		},
	})
}
