package sqlite

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var _ gokudos.KudosStorage = &KudosStorage{}

type KudosStorage struct {
	storage *gorm.DB
}

func NewKudosStorage(storage *gorm.DB) *KudosStorage {
	return &KudosStorage{storage: storage}
}

func (k *KudosStorage) StoreKudos(kudos gokudos.Kudos) error {
	return nil
}

func (k *KudosStorage) GetAllKudos() ([]gokudos.Kudos, error) {
	panic("implement me")
}

func (k *KudosStorage) ClearKudos() error {
	panic("implement me")
}
