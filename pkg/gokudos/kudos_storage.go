package gokudos

import "time"

type KudosStorage interface {
	StoreKudos(kudos Kudos) error
	GetAllKudos() ([]Kudos, error)
	GetKudosByUser(user string) ([]Kudos, error)
	DeleteKudos(message string) error
	ClearKudos() error
	SetSchedule(time time.Time) error
	ClearSchedule() error
}

type Kudos struct {
	Message     string
	SubmittedBy string
}
