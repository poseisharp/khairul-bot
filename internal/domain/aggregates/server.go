package aggregates

import "gorm.io/gorm"

type Server struct {
	gorm.Model

	ID string `gorm:"primaryKey"`
}
