package repositories

import (
	"github.com/poseisharp/khairul-bot/internal/domain/aggregates"
	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{
		db: db,
	}
}

func (r *ServerRepository) Delete(id string) error {
	return r.db.Delete(&aggregates.Server{}, id).Error
}

func (r *ServerRepository) FindAll() ([]aggregates.Server, error) {
	var servers []aggregates.Server

	if err := r.db.Find(&servers).Error; err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *ServerRepository) FindOne(id string) (*aggregates.Server, error) {
	var server aggregates.Server

	if err := r.db.First(&server, id).Error; err != nil {
		return nil, err
	}

	return &server, nil
}

func (r *ServerRepository) Store(server aggregates.Server) error {
	return r.db.Create(&server).Error
}

func (r *ServerRepository) Update(server aggregates.Server) error {
	return r.db.Save(&server).Error
}
