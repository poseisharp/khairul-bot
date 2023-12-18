package feature_jadwal

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/app/services"
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "add",
					Description: "Add jadwal preset",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "remove",
					Description: "Remove jadwal preset",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "preset",
							Description:  "Preset jadwal yang akan dihapus",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
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
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if i.ApplicationCommandData().Name == p.discordCommand.Name {
			return p.HandleCommand(s, i)
		}
	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()

		if strings.HasPrefix(data.CustomID, "jadwal-preset-add") {
			return p.handleModalAdd(s, i, data)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		return p.handleAutocomplete(s, i)
	}

	return nil
}

func (p *JadwalPresetCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	optionMap := value_objects.ArrApplicationCommandInteractionDataOption(i.ApplicationCommandData().Options).ToMap()

	mode := optionMap["mode"].StringValue()

	if mode == "add" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "jadwal-preset-add",
				Title:    "Add Jadwal Preset",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Name",
								CustomID:    "jadwal-preset-add-name",
								Required:    true,
								Placeholder: "Name",
								Style:       discordgo.TextInputShort,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Timezone",
								CustomID:    "jadwal-preset-add-timezone",
								Required:    true,
								Placeholder: "Timezone",
								Style:       discordgo.TextInputShort,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								Label:       "Latitude & Longitude",
								CustomID:    "jadwal-preset-add-lat-long",
								Required:    true,
								Placeholder: "lat,lang",
								Style:       discordgo.TextInputShort,
							},
						},
					},
				},
			},
		})

	} else if mode == "remove" {
		server, err := p.serverService.GetServer(i.GuildID)
		if err != nil {
			return err
		}

		presetId, err := strconv.Atoi(optionMap["preset"].StringValue())

		preset, err := p.presetService.GetPreset(presetId)
		if err != nil {
			return err
		}

		if preset.ServerID != server.ID {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Preset tidak ditemukan",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

	} else if mode == "list" {
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

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "List Jadwal Preset",
				Embeds:  choices,
			},
		})
	}

	return nil
}

func (p *JadwalPresetCommand) handleModalAdd(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ModalSubmitInteractionData) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Data Akan Coba di Tambahkan",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	name := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timezone := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	latLong := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	preset := entities.Preset{
		ServerID: i.GuildID,
		Name:     name,
		TimeZone: value_objects.TimeZone(timezone),
		LatLong:  value_objects.LatLong(strings.Split(latLong, ",")),
	}

	err = p.presetService.CreatePreset(preset)
	if err != nil {
		return err
	}

	if _, err := s.ChannelMessageSend(i.ChannelID, "Berhasil ditambahkan"); err != nil {
		return err
	}

	return nil
}

func (p *JadwalPresetCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
