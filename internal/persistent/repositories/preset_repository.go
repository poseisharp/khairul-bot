package repositories

import (
	"github.com/poseisharp/khairul-bot/internal/domain/aggregates"
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
	return r.db.Delete(&aggregates.Preset{}, id).Error
}

func (r *PresetRepository) FindAll() ([]aggregates.Preset, error) {
	var presets []aggregates.Preset

	if err := r.db.Find(&presets).Error; err != nil {
		return nil, err
	}

	return presets, nil
}

func (r *PresetRepository) FindOne(id int) (*aggregates.Preset, error) {
	var preset aggregates.Preset

	if err := r.db.First(&preset, id).Error; err != nil {
		return nil, err
	}

	return &preset, nil
}

func (r *PresetRepository) Store(preset aggregates.Preset) error {
	return r.db.Create(&preset).Error
}

func (r *PresetRepository) Update(preset aggregates.Preset) error {
	return r.db.Save(&preset).Error
}

func (r *PresetRepository) FindByServerID(serverID string) ([]aggregates.Preset, error) {
	var presets []aggregates.Preset

	if err := r.db.Where("server_id = ?", serverID).Find(&presets).Error; err != nil {
		return nil, err
	}

	return presets, nil
}

func (r *PresetRepository) FindByServerIDAndName(serverID string, name string) (*aggregates.Preset, error) {
	var preset aggregates.Preset

	if err := r.db.Where("server_id = ? AND name = ?", serverID, name).First(&preset).Error; err != nil {
		return nil, err
	}

	return &preset, nil
}
