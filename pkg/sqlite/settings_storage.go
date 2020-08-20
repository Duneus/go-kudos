package sqlite

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/jinzhu/gorm"
)

var _ gokudos.SettingsStorage = &SettingsStorage{}

type SettingsStorage struct {
	storage *gorm.DB
}

func NewSettingsStorage(storage *gorm.DB) *SettingsStorage {
	return &SettingsStorage{storage: storage}
}

func (s *SettingsStorage) SetScheduleSettings(scheduleSettings gokudos.Settings) error {
	model := mapSettingsToModel(scheduleSettings)

	if s.storage.NewRecord(model) == false {
		return s.storage.Create(&model).Error
	}

	var settings2 settings

	return s.storage.Find(&settings2).Where("team_id = ?", scheduleSettings.TeamId).Update("channel_id", scheduleSettings.ChannelId).Error
}

func (s *SettingsStorage) GetScheduleSettingsForTeam(teamId string) (*gokudos.Settings, error) {
	var settings2 settings
	err := s.storage.Find(&settings2).Where("team_id = ?", teamId).Error
	if err != nil {
		return nil, err
	}
	model := settings2.toModel()

	return &model, nil
}
