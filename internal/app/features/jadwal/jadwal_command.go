package feature_jadwal

import (
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/poseisharp/khairul-bot/internal/app/services"
	"github.com/poseisharp/khairul-bot/internal/domain/value_objects"
	"github.com/poseisharp/khairul-bot/internal/interfaces"
)

type JadwalCommand struct {
	interfaces.FeatureCommand

	discordCommand *discordgo.ApplicationCommand

	prayerService *services.PrayerService
	serverService *services.ServerService
	presetService *services.PresetService
}

func NewJadwalCommand(prayerService *services.PrayerService, serverService *services.ServerService, presetService *services.PresetService) *JadwalCommand {
	return &JadwalCommand{
		discordCommand: &discordgo.ApplicationCommand{
			Name:        "jadwal",
			Description: "Informasi tentang jadwal",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Description:  "Preset jadwal",
					Name:         "preset",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		prayerService: prayerService,
		serverService: serverService,
		presetService: presetService,
	}
}

func (p *JadwalCommand) DiscordCommand() *discordgo.ApplicationCommand {
	return p.discordCommand
}

func (p *JadwalCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		return p.handleAutocomplete(s, i)
	} else if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == p.discordCommand.Name {
			return p.HandleCommand(s, i)
		}
	}

	return nil
}

func (p *JadwalCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	log.Println("Handling jadwal command...")

	optionMap := value_objects.ArrApplicationCommandInteractionDataOption(i.ApplicationCommandData().Options).ToMap()

	presetName := optionMap["preset"].StringValue()

	server, err := p.serverService.GetServer(i.GuildID)
	if err != nil {
		return err
	}

	preset, err := p.presetService.GetPresetByServerIDAndName(server.ID, presetName)
	if err != nil {
		return err
	}

	schedule := p.prayerService.Calculate(preset.TimeZone, preset.LatLong)
	index := time.Now().Day() - 1

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Jadwal Sholat",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Latitude & Longitude",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Latitude",
							Value:  strconv.FormatFloat(preset.LatLong.Latitude(), 'f', 6, 32),
							Inline: true,
						},
						{
							Name:   "Longitude",
							Value:  strconv.FormatFloat(preset.LatLong.Longitude(), 'f', 6, 32),
							Inline: true,
						},
					},
				},
				{
					Title: "Jadwal Sholat",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Subuh",
							Value: schedule[index].Fajr.Format("15:04 MST"),
						},
						{
							Name:  "Dzuhur",
							Value: schedule[index].Zuhr.Format("15:04 MST"),
						},
						{
							Name:  "Ashar",
							Value: schedule[index].Asr.Format("15:04 MST"),
						},
						{
							Name:  "Maghrib",
							Value: schedule[index-1].Maghrib.Format("15:04 MST"),
						},
						{
							Name:  "Isya",
							Value: schedule[index-1].Isha.Format("15:04 MST"),
						},
					},
				},
			},
		},
	})
}

func (p *JadwalCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
}
