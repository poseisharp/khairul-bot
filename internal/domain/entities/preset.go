package entities

import (
	"github.com/poseisharp/khairul-bot/internal/domain/value_objects"
	"gorm.io/gorm"
)

type Preset struct {
	gorm.Model

	ID       int                    `gorm:"primaryKey;autoIncrement;not null"`
	Name     string                 `gorm:"not null"`
	TimeZone value_objects.TimeZone `gorm:"embedded;not null"`
	LatLong  value_objects.LatLong  `gorm:"embedded;not null"`

	ServerID string `gorm:"not null"`
	Server   Server `gorm:"foreignKey:server_id"`
}
