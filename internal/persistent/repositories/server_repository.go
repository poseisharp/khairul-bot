package repositories

import (
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
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
	return r.db.Delete(&entities.Server{}, id).Error
}

func (r *ServerRepository) FindAll() ([]entities.Server, error) {
	var servers []entities.Server

	if err := r.db.Find(&servers).Error; err != nil {
		return nil, err
	}

	return servers, nil
}

func (r *ServerRepository) FindOne(id string) (*entities.Server, error) {
	var server entities.Server

	if err := r.db.First(&server, id).Error; err != nil {
		return nil, err
	}

	return &server, nil
}

func (r *ServerRepository) Store(server entities.Server) error {
	return r.db.Create(&server).Error
}

func (r *ServerRepository) Update(server entities.Server) error {
	return r.db.Save(&server).Error
}
