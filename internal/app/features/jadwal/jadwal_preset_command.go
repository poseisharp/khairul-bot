package feature_jadwal

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/app/services"
	"github.com/poseisharp/khairul-bot/internal/domain/aggregates"
	"github.com/poseisharp/khairul-bot/internal/domain/value_objects"
	"github.com/poseisharp/khairul-bot/internal/interfaces"
)

type JadwalPresetCommand struct {
	interfaces.FeatureCommand

	discordCommand *discordgo.ApplicationCommand

	serverService *services.ServerService
	presetService *services.PresetService
}

func NewJadwalPresetCommand(serverService *services.ServerService, presetService *services.PresetService) *JadwalPresetCommand {
	return &JadwalPresetCommand{
		discordCommand: &discordgo.ApplicationCommand{
			Name:        "jadwal-preset",
			Description: "Manage jadwal preset",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Add jadwal preset",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Name of the preset",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "timezone",
							Description: "Timezone of the preset",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "lat_long",
							Description: "Latitude & Longitude of the preset",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "Remove jadwal preset",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "preset",
							Description:  "Preset to remove",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "List jadwal preset",
				},
			},
		},
		serverService: serverService,
		presetService: presetService,
	}
}

func (p *JadwalPresetCommand) DiscordCommand() *discordgo.ApplicationCommand {
	return p.discordCommand
}

func (p *JadwalPresetCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.ApplicationCommandData().Name == p.discordCommand.Name {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			return p.HandleCommand(s, i)
		case discordgo.InteractionApplicationCommandAutocomplete:
			return p.handleAutocomplete(s, i)
		}
	}
	return nil
}

func (p *JadwalPresetCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Println("Handling jadwal preset command...")

	switch i.ApplicationCommandData().Options[0].Name {
	case "add":
		return p.handleAdd(s, i)
	case "remove":
		return p.handleRemove(s, i)
	case "list":
		return p.handleList(s, i)
	}
	return nil
}

func (p *JadwalPresetCommand) handleAdd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := value_objects.ConvertInteractionDataOptionToMap(i.ApplicationCommandData().Options[0].Options)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading...",
		},
	})
	if err != nil {
		return err
	}

	name := data["name"].StringValue()
	timezone := data["timezone"].StringValue()
	latLong := data["lat_long"].StringValue()

	preset := aggregates.Preset{
		ServerID: i.GuildID,
		Name:     name,
		TimeZone: value_objects.TimeZone(timezone),
		LatLong:  value_objects.LatLong(strings.Split(latLong, ",")),
	}

	err = p.presetService.CreatePreset(preset)
	if err != nil {
		return err
	}

	content := "Preset added"
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	}); err != nil {
		return err
	}

	return nil
}

func (p *JadwalPresetCommand) handleRemove(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := value_objects.ConvertInteractionDataOptionToMap(i.ApplicationCommandData().Options[0].Options)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading...",
		},
	})
	if err != nil {
		return err
	}

	server, err := p.serverService.GetServer(i.GuildID)
	if err != nil {
		return err
	}

	presetId, err := strconv.Atoi(data["preset"].StringValue())
	if err != nil {
		return err
	}

	preset, err := p.presetService.GetPreset(presetId)
	if err != nil {
		return err
	}

	if preset.ServerID != server.ID {
		content := "Preset not found"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return err
	}

	err = p.presetService.DeletePreset(presetId)
	if err != nil {
		return err
	}

	content := "Preset removed"
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	}); err != nil {
		return err
	}

	return nil
}

func (p *JadwalPresetCommand) handleList(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading...",
		},
	})
	server, err := p.serverService.GetServer(i.GuildID)
	if err != nil {
		return err
	}

	presets, err := p.presetService.GetPresetsByServerID(server.ID)
	if err != nil {
		return err
	}

	choices := make([]*discordgo.MessageEmbed, len(presets))
	for i, preset := range presets {
		choices[i] = &discordgo.MessageEmbed{
			Title: preset.Name,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Timezone",
					Value:  string(preset.TimeZone),
					Inline: true,
				},
				{
					Name:   "Latitude",
					Value:  string(strings.Join(preset.LatLong, ", ")),
					Inline: true,
				},
			},
		}
	}

	content := "Preset List"
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
		Embeds:  &choices,
	})

	return err
}

func (p *JadwalPresetCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData().Options[0].Options[0]

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
				Value: strconv.Itoa(preset.ID),
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
