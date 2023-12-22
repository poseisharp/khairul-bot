package aggregates

import "gorm.io/gorm"

type Reminder struct {
	gorm.Model

	ID        int    `gorm:"primaryKey;autoIncrement;not null"`
	ChannelID string `gorm:"not null"`

	PresetID int    `gorm:"not null"`
	Preset   Preset `gorm:"foreignKey:preset_id"`

	ServerID string `gorm:"not null"`
	Server   Server `gorm:"foreignKey:server_id"`

	Subuh   bool `gorm:"not null;default:false"`
	Dzuhur  bool `gorm:"not null;default:false"`
	Ashar   bool `gorm:"not null;default:false"`
	Maghrib bool `gorm:"not null;default:false"`
	Isya    bool `gorm:"not null;default:false"`
}
