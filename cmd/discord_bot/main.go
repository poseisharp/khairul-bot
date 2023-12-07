package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/life4/genesis/maps"
	"github.com/life4/genesis/slices"
	"github.com/poseisharp/khairul-bot/internal/app/features"
	feature_jadwal "github.com/poseisharp/khairul-bot/internal/app/features/jadwal"
	"github.com/poseisharp/khairul-bot/internal/app/services"
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	interface_features "github.com/poseisharp/khairul-bot/internal/interfaces"
	memory_repositories "github.com/poseisharp/khairul-bot/internal/persistent/repositories/memory"
)

var (
	s        *discordgo.Session
	GuildID  string = ""
	commands        = map[string]interface_features.FeatureCommand{}
)

func addCommand(command interface_features.FeatureCommand) {
	commands[command.DiscordCommand().Name] = command
}

func init() {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN")); err != nil {
		log.Fatal("Error creating discord session")
	}

	serverRepository := memory_repositories.NewServerRepository()
	serverService := services.NewServerService(serverRepository)

	prayerService := services.NewPrayerService()

	addCommand(features.NewPingCommand())
	addCommand(feature_jadwal.NewJadwalPresetCommand(serverService))
	addCommand(feature_jadwal.NewJadwalCommand(prayerService, serverService))
	addCommand(feature_jadwal.NewJadwalManualCommand(prayerService))

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commands[i.ApplicationCommandData().Name]; ok {
				if err := h.HandleCommand(s, i); err != nil {
					log.Panicf("Error handling '%v' command: %v", i.ApplicationCommandData().Name, err)
				}
			}
		}

		for _, command := range commands {
			if err := command.Handle(s, i); err != nil {
				log.Panicf("Error handling '%v' command: %v", i.ApplicationCommandData().Name, err)
			}
		}
	})

	s.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		serverService.CreateServer(entities.Server{
			ID:            g.ID,
			JadwalPresets: []entities.JadwalPreset{},
		})
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
}

func main() {
	if err := s.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildID, slices.MapAsync(maps.Values(commands), 0, func(v interface_features.FeatureCommand) *discordgo.ApplicationCommand {
		return v.DiscordCommand()
	}))

	if err != nil {
		log.Fatalf("Cannot add commands: %v", err)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		if err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID); err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}

// func initDb() (db *gorm.DB, err error) {
// 	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	if err != nil {
// 		return
// 	}

// 	db.AutoMigrate(&entities.Server{})

// 	return
// }
