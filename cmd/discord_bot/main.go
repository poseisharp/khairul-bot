package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/life4/genesis/maps"
	"github.com/life4/genesis/slices"
	ping_feature "github.com/poseisharp/khairul-bot/internal/app/features/ping"
	"github.com/poseisharp/khairul-bot/internal/interfaces/features"
)

var (
	s        *discordgo.Session
	GuildID  string = ""
	commands        = map[string]features.FeatureCommand{}
)

func addCommand(command features.FeatureCommand) {
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

	addCommand(ping_feature.New())

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands[i.ApplicationCommandData().Name]; ok {
			if err := h.Handle(s, i); err != nil {
				log.Fatalf("Error handling '%v' command: %v", i.ApplicationCommandData().Name, err)
			}
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err := s.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	log.Println("Adding commands...")
	registeredCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildID, slices.MapAsync(maps.Values(commands), 0, func(v features.FeatureCommand) *discordgo.ApplicationCommand { return v.DiscordCommand() }))

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
