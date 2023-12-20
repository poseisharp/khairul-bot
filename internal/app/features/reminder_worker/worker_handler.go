package feature_reminder_worker

import (
	"log"
	"time"

	"github.com/gocraft/work"
	"github.com/poseisharp/khairul-bot/internal/app/services"
)

type ReminderWorkerHandler struct {
	enqueuer *work.Enqueuer

	reminderService *services.ReminderService
	prayerService   *services.PrayerService
	discordService  *services.DiscordService
}

func NewReminderWorkerHandler(enqueuer *work.Enqueuer, reminderService *services.ReminderService, prayerService *services.PrayerService, discordService *services.DiscordService) *ReminderWorkerHandler {
	return &ReminderWorkerHandler{
		enqueuer:        enqueuer,
		reminderService: reminderService,
		prayerService:   prayerService,
		discordService:  discordService,
	}
}

func (h *ReminderWorkerHandler) SetupReminder(job *work.Job) error {
	reminders, err := h.reminderService.GetReminders()
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		now := time.Now().UTC()

		schedules := h.prayerService.Calculate(reminder.Preset.TimeZone, reminder.Preset.LatLong)
		dayOfYear := now.YearDay() - 1

		if reminder.Subuh {
			h.enqueuer.EnqueueIn("run_reminder", int64(time.Until(schedules[dayOfYear].Fajr).Seconds()), work.Q{
				"reminder_id": reminder.ID,
				"prayer":      "subuh",
				"schedule":    schedules[dayOfYear].Fajr.Format("15:04 MST"),
			})
		}
		if reminder.Dzuhur {
			h.enqueuer.EnqueueIn("run_reminder", int64(time.Until(schedules[dayOfYear].Zuhr).Seconds()), work.Q{
				"reminder_id": reminder.ID,
				"prayer":      "dzuhur",
				"schedule":    schedules[dayOfYear].Zuhr.Format("15:04 MST"),
			})
		}
		if reminder.Ashar {
			h.enqueuer.EnqueueIn("run_reminder", int64(time.Until(schedules[dayOfYear].Asr).Seconds()), work.Q{
				"reminder_id": reminder.ID,
				"prayer":      "ashar",
				"schedule":    schedules[dayOfYear].Asr.Format("15:04 MST"),
			})
		}
		if reminder.Maghrib {
			h.enqueuer.EnqueueIn("run_reminder", int64(time.Until(schedules[dayOfYear].Maghrib).Seconds()), work.Q{
				"reminder_id": reminder.ID,
				"prayer":      "maghrib",
				"schedule":    schedules[dayOfYear-1].Maghrib.Format("15:04 MST"),
			})
		}
		if reminder.Isya {
			h.enqueuer.EnqueueIn("run_reminder", int64(time.Until(schedules[dayOfYear].Isha).Seconds()), work.Q{
				"reminder_id": reminder.ID,
				"prayer":      "isya",
				"schedule":    schedules[dayOfYear-1].Isha.Format("15:04 MST"),
			})
		}
	}

	return nil
}

func (h *ReminderWorkerHandler) RunReminder(job *work.Job) error {
	reminderId := job.ArgInt64("reminder_id")
	prayer := job.ArgString("prayer")
	schedule := job.ArgString("schedule")

	log.Println("Handle run reminder")

	reminder, err := h.reminderService.GetReminder(int(reminderId))
	if err != nil {
		log.Println(err)
		return err
	}

	err = h.discordService.SendTextMessage(reminder.ChannelID, "Waktunya sholat "+prayer+" ("+schedule+")")

	if err != nil {
		log.Println(err)
	}

	return err
}
