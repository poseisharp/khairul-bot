package value_objects

import "github.com/bwmarrin/discordgo"

type ArrApplicationCommandInteractionDataOption []*discordgo.ApplicationCommandInteractionDataOption
type MapApplicationCommandInteractionDataOption map[string]*discordgo.ApplicationCommandInteractionDataOption

func (a ArrApplicationCommandInteractionDataOption) ToMap() MapApplicationCommandInteractionDataOption {
	m := make(MapApplicationCommandInteractionDataOption)
	for _, v := range a {
		m[v.Name] = v
	}
	return m
}
