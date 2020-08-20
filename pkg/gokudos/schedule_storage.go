package gokudos

type ScheduleStorage interface {
	GetSchedule(settings Settings) (Schedule, error)
	SetSchedule(schedule Schedule) error
	ClearSchedule(schedule Schedule) error
}

type Schedule struct {
	TeamId      string
	ChannelId   string
	ScheduleId  string
	ScheduledAt int64
}
