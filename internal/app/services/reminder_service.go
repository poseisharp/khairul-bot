package services

import (
	"github.com/poseisharp/khairul-bot/internal/domain/aggregates"
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

func (s *ReminderService) GetReminder(id int) (*aggregates.Reminder, error) {
	reminder, err := s.reminderRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (s *ReminderService) GetReminders() ([]aggregates.Reminder, error) {
	reminders, err := s.reminderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (s *ReminderService) CreateReminder(reminder aggregates.Reminder) error {
	return s.reminderRepository.Store(reminder)
}

func (s *ReminderService) UpdateReminder(reminder aggregates.Reminder) error {
	return s.reminderRepository.Update(reminder)
}

func (s *ReminderService) DeleteReminder(id int) error {
	return s.reminderRepository.Delete(id)
}

func (s *ReminderService) GetRemindersByServerID(serverID string) ([]aggregates.Reminder, error) {
	reminders, err := s.reminderRepository.FindByServerID(serverID)
	if err != nil {
		return nil, err
	}

	return reminders, nil
}
