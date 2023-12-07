package services

import (
	"time"

	"github.com/hablullah/go-prayer"
	"github.com/poseisharp/khairul-bot/internal/domain/entities"
)

type PrayerService struct {
}

func NewPrayerService() *PrayerService {
	return &PrayerService{}
}

func (s *PrayerService) Calculate(timezone entities.TimeZone, latLong entities.LatLong) []prayer.Schedule {
	schedule, _ := prayer.Calculate(prayer.Config{
		Latitude:           latLong.Latitude(),
		Longitude:          latLong.Longitude(),
		Timezone:           timezone.LoadLocation(),
		TwilightConvention: prayer.Kemenag(),
		AsrConvention:      prayer.Shafii,
		PreciseToSeconds:   true,
		Corrections:        prayer.ScheduleCorrections{},
	}, time.Now().Year())

	return schedule
}
