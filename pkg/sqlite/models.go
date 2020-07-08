package sqlite

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
)

type kudos struct {
	ID          int `gorm:"primary_key"`
	Message     string
	SubmittedBy string
	SubmittedIn string
}

type settings struct {
	TeamId      string `gorm:"unique;not null"`
	ChannelId   string
	ScheduledId string
}

func mapKudosToModel(kudos2 gokudos.Kudos) *kudos {
	return &kudos{
		ID:          kudos2.ID,
		Message:     kudos2.Message,
		SubmittedBy: kudos2.SubmittedBy,
		SubmittedIn: kudos2.SubmittedIn,
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

func mapSettingsToModel(settings2 gokudos.Settings) *settings {
	return &settings{
		TeamId:      settings2.TeamId,
		ChannelId:   settings2.ChannelId,
		ScheduledId: "",
	}
}

func (s *settings) toModel() gokudos.Settings {
	return gokudos.Settings{
		TeamId: s.TeamId,
		ChannelId: s.ChannelId,

	}
}
