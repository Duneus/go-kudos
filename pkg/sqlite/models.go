package sqlite

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"time"
)

type kudos struct {
	Message     string
	SubmittedBy string
}

type schedule struct {
	Scheduled time.Time `gorm:"not null"`
}

func mapKudosToModel(kudos2 gokudos.Kudos) *kudos {
	return &kudos{
		Message: kudos2.Message,
		SubmittedBy: kudos2.SubmittedBy,
	}
}

func (k *kudos) toModel() gokudos.Kudos {
	return gokudos.Kudos{
		Message:     k.Message,
		SubmittedBy: k.SubmittedBy,
	}
}
