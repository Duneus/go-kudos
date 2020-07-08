package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Duneus/go-kudos/pkg/config"
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
	"strconv"
	"time"
)

var _ gokudos.KudosService = &KudosService{}

type KudosService struct {
	kudosStorage    gokudos.KudosStorage
	scheduleStorage gokudos.ScheduleStorage
	settingsStorage gokudos.SettingsStorage
	cfg             config.Config
	client          *slack.Client
}

func NewKudosService(
	kudosStorage gokudos.KudosStorage,
	scheduleStorage gokudos.ScheduleStorage,
	settingsStorage gokudos.SettingsStorage,
	cfg config.Config,
	client *slack.Client,
) *KudosService {
	return &KudosService{
		kudosStorage:    kudosStorage,
		scheduleStorage: scheduleStorage,
		settingsStorage: settingsStorage,
		cfg:             cfg,
		client:          client,
	}
}

func (s *KudosService) Forward(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, err := s.parseMessage(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.handleEvent(eventsAPIEvent, body, w)
}

func (s *KudosService) AddKudos(rw http.ResponseWriter, r *http.Request) {
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		panic(err)
	}

	kudos := gokudos.Kudos{
		Message:     command.Text,
		SubmittedBy: command.UserID,
		SubmittedIn: command.TeamID,
	}

	s.kudosStorage.StoreKudos(kudos)

	s.client.PostEphemeral(command.ChannelID, command.UserID, slack.MsgOptionText("Thanks for sending your kudos!", false))

	rw.WriteHeader(200)
}

func (s *KudosService) PublishKudos(rw http.ResponseWriter, r *http.Request) {
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		panic(err)
	}
	kudos, _ := s.kudosStorage.GetAllKudosInTeam(command.TeamID)

	var message string

	for _, k := range kudos {
		message = message + prependHeart(k.Message) + "\n"
	}

	s.client.PostMessage(command.ChannelID, slack.MsgOptionText(message, false))

	rw.WriteHeader(200)
}

func (s *KudosService) HandleInteractivity(rw http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	for _, action := range payload.ActionCallback.BlockActions {
		var view slack.HomeTabViewRequest
		fmt.Printf("Action idx is: %s\n", action.ActionID)
		switch action.ActionID {
		case "show_kudos":
			kudos, err := s.kudosStorage.GetKudosByUser(payload.User.ID)
			if err != nil {
				fmt.Printf("error: %w", err)
			}
			view, err = handleAppHomeTabWithKudosList(kudos)
			if err != nil {
				fmt.Printf("error: %w", err)
			}
			_, err = s.client.PublishView(payload.User.ID, view, "")
			if err != nil {
				fmt.Printf("error: %w", err)
			}
		case "show_all_kudos":
			kudos, err := s.kudosStorage.GetAllKudosInTeam(payload.Team.ID)
			if err != nil {
				fmt.Printf("error: %w", err)
			}
			view, err = handleAppHomeTabWithKudosList(kudos)
			if err != nil {
				fmt.Printf("error: %w", err)
			}
			s.client.PublishView(payload.User.ID, view, "")
		case "remove_kudos":
			kudosId, err := strconv.Atoi(action.Value)
			if err != nil {
				panic(err)
			}
			s.kudosStorage.DeleteKudos(kudosId)
			kudos, _ := s.kudosStorage.GetKudosByUser(payload.User.ID)
			view, err = handleAppHomeTabWithKudosList(kudos)
			s.client.PublishView(payload.User.ID, view, "")
		case "schedule":
			view, err = handleAppHomeSchedulingTab()
			s.client.PublishView(payload.User.ID, view, "")
		case "schedule_new":
			err = s.scheduleStorage.SetSchedule(gokudos.Schedule{
				TeamId:      payload.Team.ID,
				ChannelId:   payload.Channel.ID,
				ScheduleId:  "",
				ScheduledAt: time.Now().Add(time.Minute).Unix(),
			})
			if err != nil {
				fmt.Printf("Error while setting a scheduled message: %v\n", err)
			}
		case "channel_select":
			channel := ""
			settings, err := s.settingsStorage.GetScheduleSettingsForTeam(payload.Team.ID)
			if err != nil || settings == nil {
				fmt.Printf("Print error %+v", err)
			}
			fmt.Printf("SelectedChannel: %+v\n", settings)
			if settings != nil {
				channel = settings.ChannelId
			}
			view, err = handleAppHomeTabChannelPicker(channel)
			_, err = s.client.PublishView(payload.User.ID, view, "")
			if err != nil {
				fmt.Printf("Print error %+v", err)
			}
		case "select_channel":
			settings := gokudos.Settings{
				TeamId:    payload.Team.ID,
				ChannelId: action.SelectedChannel,
			}
			s.settingsStorage.SetScheduleSettings(settings)
			view, err = handleAppHomeTabChannelPicker(settings.ChannelId)
			_, err = s.client.PublishView(payload.User.ID, view, "")
			if err != nil {
				fmt.Printf("Print error %+v", err)
			}
		case "hide_kudos":
			view, err = handleAppHomeTab()
			s.client.PublishView(payload.User.ID, view, "")
		default:
			view, err = handleAppHomeTab()
			s.client.PublishView(payload.User.ID, view, "")
		}
	}

	rw.WriteHeader(200)
}

func (s *KudosService) sendMessage(channelId string, message string) error {
	_, _, err := s.client.PostMessage(channelId, slack.MsgOptionText(message, false))

	return err
}

func (s *KudosService) parseMessage(body string) (event slackevents.EventsAPIEvent, err error) {
	return slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: s.cfg.VerificationToken}))
}

func (s *KudosService) handleEvent(event slackevents.EventsAPIEvent, body string, w http.ResponseWriter) {
	if event.Type == slackevents.URLVerification {
		s.handleVerification(body, w)
	}
	if event.Type == slackevents.CallbackEvent {
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppHomeOpenedEvent:
			view, _ := handleAppHomeTab()
			s.client.PublishView(ev.User, view, "")
		case *slackevents.AppMentionEvent:
			s.sendMessage(ev.Channel, "Yes, hello.")
		case *slackevents.MessageEvent:
			err := s.handleMessageEvents(ev)
			if err != nil {
				fmt.Printf("Error handling message: %v", err)
			}
		case *slackevents.MessageAction:
		}
	}

	w.WriteHeader(200)
}

func (s *KudosService) handleVerification(body string, w http.ResponseWriter) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text")
	w.Write([]byte(r.Challenge))
}

func (s *KudosService) handleMessageEvents(event *slackevents.MessageEvent) error {
	if event.BotID == "" {
		s.sendMessage(event.Channel, "Hello, thanks")
	}

	return nil
}
