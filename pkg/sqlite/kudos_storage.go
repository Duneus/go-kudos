package sqlite

import (
	"fmt"
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
	return k.storage.Create(mapKudosToModel(kudos)).Error
}

func (k *KudosStorage) GetAllKudosInTeam(teamId string) ([]gokudos.Kudos, error) {
	var kudos []kudos

	err := k.storage.Find(&kudos).Where("submitted_in = ?", teamId).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving all kudos: %w", err)
	}

	var allKudos []gokudos.Kudos

	for _, k := range kudos {
		allKudos = append(allKudos, k.toModel())
	}

	return allKudos, nil
}

func (k *KudosStorage) GetKudosByUser(user string) ([]gokudos.Kudos, error) {
	var kudos []kudos

	err := k.storage.Where("submitted_by = ?", user).Find(&kudos).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving kudos: %w", err)
	}

	var allKudos []gokudos.Kudos

	for _, k := range kudos {
		allKudos = append(allKudos, k.toModel())
	}

	return allKudos, nil

}

func (k *KudosStorage) DeleteKudos(kudosId int) error {
	return k.storage.Delete(kudos{}, "id = ?", kudosId).Error
}

func (k *KudosStorage) ClearKudos(teamId string) error {
	return k.storage.Delete(&kudos{}).Where("team_id = ?", teamId).Error
}

func (k *KudosStorage) SetSchedule(schedule gokudos.Schedule) error {
	return k.storage.Create(mapScheduleToModel(schedule)).Error
}

func (k *KudosStorage) ClearSchedule(teamId string) error {
	return k.storage.Delete(&schedule{}).Where("team_id = ?", teamId).Error
}
