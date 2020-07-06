package gokudos

import "time"

type KudosStorage interface {
	StoreKudos(kudos Kudos) error
	GetAllKudosInTeam(teamId string) ([]Kudos, error)
	GetKudosByUser(user string) ([]Kudos, error)
	DeleteKudos(kudosId string) error
	ClearKudos(teamId string) error
	SetSchedule(schedule Schedule) error
	ClearSchedule(teamId string) error
}

type Kudos struct {
	ID          string
	Message     string
	SubmittedBy string
	SubmittedIn string
}

type Schedule struct {
	TeamId string
	Time   time.Time
}
