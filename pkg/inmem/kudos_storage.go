package inmem

import "github.com/Duneus/go-kudos/pkg/gokudos"

var _ gokudos.KudosStorage = &KudosStorage{}

type KudosStorage struct {
	storage []gokudos.Kudos
}

func NewKudosStorage() *KudosStorage {
	return &KudosStorage{}
}

func (k *KudosStorage) StoreKudos(kudos gokudos.Kudos) error {
	k.storage = append(k.storage, kudos)
	return nil
}

func (k *KudosStorage) GetAllKudos() ([]gokudos.Kudos, error) {
	list := make([]gokudos.Kudos, len(k.storage))
	copy(list, k.storage)

	return list, nil
}

func (k *KudosStorage) ClearKudos() error {
	k.storage = nil
	return nil
}

