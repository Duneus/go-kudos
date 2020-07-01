package inmem

import "github.com/Duneus/go-kudos/pkg/gokudos"

var _ gokudos.KudosStorage = &KudosStorage{}

type KudosStorage struct {}

func NewKudosStorage() *KudosStorage {
	return &KudosStorage{}
}

func (k *KudosStorage) StoreKudos(kudos gokudos.Kudos) error {
	panic("implement me")
}

func (k *KudosStorage) GetAllKudos() ([]gokudos.Kudos, error) {
	panic("implement me")
}

func (k *KudosStorage) ClearKudos() error {
	panic("implement me")
}

