package slack

import (
	"fmt"
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/slack-go/slack"
	"strconv"
)

var _ gokudos.ScheduleStorage = &ScheduleStorage{}

func NewScheduleStorage(client *slack.Client) *ScheduleStorage {
	return &ScheduleStorage{client: client}
}

type ScheduleStorage struct {
	client *slack.Client
}

func (s *ScheduleStorage) GetSchedule(settings gokudos.Settings) (gokudos.Schedule, error) {
	messages, _, _ := s.client.GetScheduledMessages(&slack.GetScheduledMessagesParameters{Channel: settings.ChannelId})

	fmt.Printf("GetSchedule returned: %+v\n", messages)
	schedule := gokudos.Schedule{
		TeamId:      "message.Team",
		ChannelId:   "message.Channel",
		ScheduleId:  "message",
		ScheduledAt: 0,
	}

	return schedule, nil
}

func (s *ScheduleStorage) SetScheduleSettings(settings gokudos.Settings) error {
	panic("")
}

func (s *ScheduleStorage) SetSchedule(schedule gokudos.Schedule) error {
	scheduledAt := strconv.FormatInt(schedule.ScheduledAt, 10)
	_, _, err := s.client.ScheduleMessage(schedule.ChannelId, scheduledAt)
	if err != nil {
		return fmt.Errorf("zajebisty kurwa error: %w", err)
	}

	return nil
}

func (s *ScheduleStorage) ClearSchedule(schedule gokudos.Schedule) error {
	_, err := s.client.DeleteScheduledMessage(&slack.DeleteScheduledMessageParameters{
		Channel:            schedule.ChannelId,
		ScheduledMessageID: schedule.ScheduleId,
		AsUser:             false,
	})

	return err
}
