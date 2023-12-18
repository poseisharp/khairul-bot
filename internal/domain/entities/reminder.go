package entities

type Reminder struct {
	ID        string `gorm:"primaryKey;autoIncrement;not null"`
	ChannelID string `gorm:"not null"`

	PresetID int    `gorm:"not null"`
	Preset   Preset `gorm:"foreignKey:preset_id"`

	ServerID string `gorm:"not null"`
	Server   Server `gorm:"foreignKey:server_id"`
}
