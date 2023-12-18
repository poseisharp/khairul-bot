package repositories

import (
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	"gorm.io/gorm"
)

type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) *ReminderRepository {
	return &ReminderRepository{
		db: db,
	}
}

func (r *ReminderRepository) Delete(id string) error {
	return r.db.Delete(&entities.Reminder{}, id).Error
}

func (r *ReminderRepository) FindAll() ([]entities.Reminder, error) {
	var reminders []entities.Reminder

	if err := r.db.Find(&reminders).Error; err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *ReminderRepository) FindOne(id string) (*entities.Reminder, error) {
	var reminder entities.Reminder

	if err := r.db.First(&reminder, id).Error; err != nil {
		return nil, err
	}

	return &reminder, nil
}

func (r *ReminderRepository) Store(reminder entities.Reminder) error {
	return r.db.Create(&reminder).Error
}

func (r *ReminderRepository) Update(reminder entities.Reminder) error {
	return r.db.Save(&reminder).Error
}

func (r *ReminderRepository) FindByServerID(serverID string) ([]entities.Reminder, error) {
	var reminders []entities.Reminder

	if err := r.db.Where("server_id = ?", serverID).Find(&reminders).Error; err != nil {
		return nil, err
	}

	return reminders, nil
}
