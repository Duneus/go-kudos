package gokudos

type KudosStorage interface {
	StoreKudos(kudos Kudos) error
	GetAllKudos() ([]Kudos, error)
	ClearKudos() error
}

type Kudos struct {
	Message string
}
