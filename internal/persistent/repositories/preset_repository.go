package repositories

import (
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	"gorm.io/gorm"
)

type PresetRepository struct {
	db *gorm.DB
}

func NewPresetRepository(db *gorm.DB) *PresetRepository {
	return &PresetRepository{
		db: db,
	}
}

func (r *PresetRepository) Delete(id int) error {
	return r.db.Delete(&entities.Preset{}, id).Error
}

func (r *PresetRepository) FindAll() ([]entities.Preset, error) {
	var presets []entities.Preset

	if err := r.db.Find(&presets).Error; err != nil {
		return nil, err
	}

	return presets, nil
}

func (r *PresetRepository) FindOne(id int) (*entities.Preset, error) {
	var preset entities.Preset

	if err := r.db.First(&preset, id).Error; err != nil {
		return nil, err
	}

	return &preset, nil
}

func (r *PresetRepository) Store(preset entities.Preset) error {
	return r.db.Create(&preset).Error
}

func (r *PresetRepository) Update(preset entities.Preset) error {
	return r.db.Save(&preset).Error
}

func (r *PresetRepository) FindByServerID(serverID string) ([]entities.Preset, error) {
	var presets []entities.Preset

	if err := r.db.Where("server_id = ?", serverID).Find(&presets).Error; err != nil {
		return nil, err
	}

	return presets, nil
}

func (r *PresetRepository) FindByServerIDAndName(serverID string, name string) (*entities.Preset, error) {
	var preset entities.Preset

	if err := r.db.Where("server_id = ? AND name = ?", serverID, name).First(&preset).Error; err != nil {
		return nil, err
	}

	return &preset, nil
}
