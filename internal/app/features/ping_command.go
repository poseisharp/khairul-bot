package features

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/interfaces"
)

type PingCommand struct {
	interfaces.FeatureCommand

	discordCommand *discordgo.ApplicationCommand
}

func NewPingCommand() *PingCommand {
	return &PingCommand{
		discordCommand: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Ping the bot",
		},
	}
}

func (p *PingCommand) DiscordCommand() *discordgo.ApplicationCommand {
	return p.discordCommand
}

func (*PingCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Println("Handling ping command...")
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}

func (*PingCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return nil
}
