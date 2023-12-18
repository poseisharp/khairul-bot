package services

import (
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
	"github.com/poseisharp/khairul-bot/internal/persistent/repositories"
)

type ReminderService struct {
	reminderRepository *repositories.ReminderRepository
}

func NewReminderService(reminderRepository *repositories.ReminderRepository) *ReminderService {
	return &ReminderService{
		reminderRepository: reminderRepository,
	}
}

func (s *ReminderService) GetReminder(id string) (*entities.Reminder, error) {
	reminder, err := s.reminderRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (s *ReminderService) GetReminders() ([]entities.Reminder, error) {
	reminders, err := s.reminderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (s *ReminderService) CreateReminder(reminder entities.Reminder) error {
	return s.reminderRepository.Store(reminder)
}

func (s *ReminderService) UpdateReminder(reminder entities.Reminder) error {
	return s.reminderRepository.Update(reminder)
}

func (s *ReminderService) DeleteReminder(id string) error {
	return s.reminderRepository.Delete(id)
}

func (s *ReminderService) GetRemindersByServerID(serverID string) ([]entities.Reminder, error) {
	reminders, err := s.reminderRepository.FindByServerID(serverID)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}
