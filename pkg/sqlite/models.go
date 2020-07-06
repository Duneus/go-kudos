package sqlite

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"time"
)

type kudos struct {
	ID          string
	Message     string
	SubmittedBy string
	SubmittedIn string
}

type schedule struct {
	TeamId string    `gorm:"unique;not null"`
	Time   time.Time `gorm:"not null"`
}

func mapKudosToModel(kudos2 gokudos.Kudos) *kudos {
	return &kudos{
		ID:          kudos2.ID,
		Message:     kudos2.Message,
		SubmittedBy: kudos2.SubmittedBy,
		SubmittedIn: kudos2.SubmittedIn,
	}
}

func mapScheduleToModel(schedule2 gokudos.Schedule) *schedule {
	return &schedule{
		TeamId: schedule2.TeamId,
		Time:   schedule2.Time,
	}
}

func (k *kudos) toModel() gokudos.Kudos {
	return gokudos.Kudos{
		ID:          k.ID,
		Message:     k.Message,
		SubmittedBy: k.SubmittedBy,
		SubmittedIn: k.SubmittedIn,
	}
}
