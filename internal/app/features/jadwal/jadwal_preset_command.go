package feature_jadwal

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/app/services"
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	"github.com/poseisharp/khairul-bot/internal/interfaces"
)

type JadwalPresetCommand struct {
	interfaces.FeatureCommand

	discordCommand *discordgo.ApplicationCommand

	serverService *services.ServerService
}

func NewJadwalPresetCommand(serverService *services.ServerService) *JadwalPresetCommand {
	return &JadwalPresetCommand{
		discordCommand: &discordgo.ApplicationCommand{
			Name:        "jadwal-preset",
			Description: "Manage jadwal preset",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mode",
					Description: "Mode",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "add",
							Value: "add",
						},
						{
							Name:  "remove",
							Value: "remove",
						},
						{
							Name:  "list",
							Value: "list",
						},
					},
				},
			},
		},
		serverService: serverService,
	}
}

func (p *JadwalPresetCommand) DiscordCommand() *discordgo.ApplicationCommand {
	return p.discordCommand
}

func (p *JadwalPresetCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

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

		presets := make([]discordgo.SelectMenuOption, len(server.JadwalPresets))
		for i, preset := range server.JadwalPresets {
			presets[i] = discordgo.SelectMenuOption{
				Label: preset.Name,
				Value: preset.Name,
			}
		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "jadwal-preset-remove",
				Title:    "Remove Jadwal Preset",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID: "jadwal-preset-remove-select",
								Options:  presets,
								MenuType: discordgo.StringSelectMenu,
							},
						},
					},
				},
			},
		})
	} else if mode == "list" {
		server, err := p.serverService.GetServer(i.GuildID)
		if err != nil {
			return err
		}

		presets := make([]*discordgo.MessageEmbed, len(server.JadwalPresets))
		for i, preset := range server.JadwalPresets {
			presets[i] = &discordgo.MessageEmbed{
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
				Embeds:  presets,
			},
		})
	}

	return nil
}

func (p *JadwalPresetCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type == discordgo.InteractionModalSubmit {
		data := i.ModalSubmitData()

		if strings.HasPrefix(data.CustomID, "jadwal-preset-add") {
			p.handleModalAdd(s, i, data)
		} else if strings.HasPrefix(data.CustomID, "jadwal-preset-remove") {
			p.handleModalRemove(s, i, data)
		}

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

	preset := entities.JadwalPreset{
		Name:     name,
		TimeZone: entities.TimeZone(timezone),
		LatLong:  entities.LatLong(strings.Split(latLong, ",")),
	}

	server, err := p.serverService.GetServer(i.GuildID)

	if err != nil {
		return err
	}

	server.JadwalPresets = append(server.JadwalPresets, preset)

	if err := p.serverService.UpdateServer(*server); err != nil {
		return err
	}

	if _, err := s.ChannelMessageSend(i.ChannelID, "Berhasil ditambahkan"); err != nil {
		return err
	}

	return nil
}

func (p *JadwalPresetCommand) handleModalRemove(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ModalSubmitInteractionData) error {
	preset, err := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu).MarshalJSON()

	// server, err := p.serverService.GetServer(i.GuildID)

	if err != nil {
		return err
	}

	// for i, preset := range server.JadwalPresets {
	// 	if preset.Name == selectMenu.Value {
	// 		server.JadwalPresets = append(server.JadwalPresets[:i], server.JadwalPresets[i+1:]...)
	// 		break
	// 	}
	// }

	// if err := p.serverService.UpdateServer(*server); err != nil {
	// 	return err
	// }

	if _, err := s.ChannelMessageSend(i.ChannelID, string(preset)); err != nil {
		return err
	}

	return nil
}
