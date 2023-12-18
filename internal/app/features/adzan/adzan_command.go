package feature_adzan

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/app/services"
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	"github.com/poseisharp/khairul-bot/internal/domain/value_objects"
	"github.com/poseisharp/khairul-bot/internal/interfaces"
)

type AdzanCommand struct {
	interfaces.FeatureCommand

	discordCommand *discordgo.ApplicationCommand

	serverService   *services.ServerService
	presetService   *services.PresetService
	reminderService *services.ReminderService
}

func NewAdzanCommand(serverService *services.ServerService, reminderService *services.ReminderService) *AdzanCommand {
	return &AdzanCommand{
		discordCommand: &discordgo.ApplicationCommand{
			Name:        "adzan",
			Description: "Mengatur reminder adzan",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "adzan-set",
					Description: "Mengatur reminder adzan",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "preset",
							Description: "Preset reminder adzan",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "channel",
							Description: "Channel untuk reminder adzan",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "adzan-unset",
					Description: "Menghapus reminder adzan",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "reminder",
							Description:  "Reminder adzan yang akan dihapus",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
			},
		},
		serverService:   serverService,
		reminderService: reminderService,
	}
}

func (p *AdzanCommand) DiscordCommand() *discordgo.ApplicationCommand {
	return p.discordCommand
}

func (p *AdzanCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == p.discordCommand.Name {
			return p.HandleCommand(s, i)
		}
	} else if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		return p.handleAutocomplete(s, i)
	}

	return nil
}

func (p *AdzanCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Println("Handling adzan command...")

	optionMap := value_objects.ArrApplicationCommandInteractionDataOption(i.ApplicationCommandData().Options).ToMap()

	if i.ApplicationCommandData().Options[0].Name == "adzan-set" {
		return p.handleSet(s, i, optionMap)
	} else if i.ApplicationCommandData().Options[0].Name == "adzan-unset" {
		return p.handleUnset(s, i, optionMap)
	}

	return nil
}

func (p *AdzanCommand) handleSet(s *discordgo.Session, i *discordgo.InteractionCreate, optionMap value_objects.MapApplicationCommandInteractionDataOption) error {
	serverID := i.GuildID
	presetName := optionMap["preset"].StringValue()
	channelID := optionMap["channel"].StringValue()

	preset, err := p.presetService.GetPresetByServerIDAndName(serverID, presetName)
	if err != nil {
		return err
	}

	err = p.reminderService.CreateReminder(entities.Reminder{
		ServerID:  serverID,
		PresetID:  preset.ID,
		ChannelID: channelID,
	})

	return nil
}

func (p *AdzanCommand) handleUnset(s *discordgo.Session, i *discordgo.InteractionCreate, optionMap value_objects.MapApplicationCommandInteractionDataOption) error {
	reminder := optionMap["reminder"].StringValue()

	err := p.reminderService.DeleteReminder(reminder)
	if err != nil {
		return err
	}

	return nil
}

func (p *AdzanCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()

	if data.Name == "preset" {
		server, err := p.serverService.GetServer(i.GuildID)
		if err != nil {
			return err
		}

		presets, err := p.presetService.GetPresetsByServerID(server.ID)
		if err != nil {
			return err
		}

		choices := make([]*discordgo.ApplicationCommandOptionChoice, len(presets))
		for i, preset := range presets {
			choices[i] = &discordgo.ApplicationCommandOptionChoice{
				Name:  preset.Name,
				Value: preset.Name,
			}
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	} else if data.Name == "reminder" {
		reminders, err := p.reminderService.GetRemindersByServerID(i.GuildID)
		if err != nil {
			return err
		}

		choices := make([]*discordgo.ApplicationCommandOptionChoice, len(reminders))
		for i, reminder := range reminders {
			choices[i] = &discordgo.ApplicationCommandOptionChoice{
				Name:  reminder.ID,
				Value: reminder.ID,
			}
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	}

	return nil
}
