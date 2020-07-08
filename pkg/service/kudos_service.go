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

const (
	ShowKudosView         = "show_kudos_view"
	ShowAllKudosView      = "show_all_kudos_view"
	ShowSchedulingView    = "show_scheduling_view"
	ShowChannelSelectView = "show_channel_select_view"

	SetSchedule = "set_schedule"
	SetChannel  = "set_channel"
	RemoveKudos = "remove_kudos"
)

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
	_, _ = buf.ReadFrom(r.Body)
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

	err = s.kudosStorage.StoreKudos(kudos)
	if err != nil {
		fmt.Printf("Print error %+v", err)
	}

	_, err = s.client.PostEphemeral(command.ChannelID, command.UserID, slack.MsgOptionText("Thanks for sending your kudos!", false))
	if err != nil {
		fmt.Printf("Print error %+v", err)
	}

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

	_, _, err = s.client.PostMessage(command.ChannelID, slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("Print error %+v", err)
	}

	rw.WriteHeader(200)
}

func (s *KudosService) HandleInteractivity(rw http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	for _, action := range payload.ActionCallback.BlockActions {
		switch action.ActionID {
		case ShowKudosView:
			s.showMyKudosView(payload)
		case ShowAllKudosView:
			s.showAllKudosView(payload)
		case RemoveKudos:
			s.removeKudos(payload, action)
		case ShowSchedulingView:
			s.showSchedulingView(payload)
		case SetSchedule:
			s.setSchedule(payload)
		case ShowChannelSelectView:
			s.showChannelSelectView(payload)
		case SetChannel:
			s.setChannel(payload, action)
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
			_, err := s.client.PublishView(ev.User, view, "")
			if err != nil {
				fmt.Printf("Print error %+v", err)
			}
		case *slackevents.AppMentionEvent:
			err := s.sendMessage(ev.Channel, "Yes, hello.")
			if err != nil {
				fmt.Printf("Print error %+v", err)
			}
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
	_, _ = w.Write([]byte(r.Challenge))
}

func (s *KudosService) handleMessageEvents(event *slackevents.MessageEvent) error {
	if event.BotID == "" {
		err := s.sendMessage(event.Channel, "Hello, thanks")
		if err != nil {
			fmt.Printf("Print error %+v", err)
		}
	}

	return nil
}

func (s *KudosService) showMyKudosView(payload slack.InteractionCallback) {
	kudos, err := s.kudosStorage.GetKudosByUser(payload.User.ID)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	view, err := handleAppHomeTabWithKudosList(kudos)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	_, err = s.client.PublishView(payload.User.ID, view, "")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func (s *KudosService) showAllKudosView(payload slack.InteractionCallback) {
	kudos, err := s.kudosStorage.GetAllKudosInTeam(payload.Team.ID)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	view, err := handleAppHomeTabWithKudosList(kudos)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	_, err = s.client.PublishView(payload.User.ID, view, "")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func (s *KudosService) showSchedulingView(payload slack.InteractionCallback) {
	view, err := handleAppHomeSchedulingTab()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	_, err = s.client.PublishView(payload.User.ID, view, "")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func (s *KudosService) showChannelSelectView(payload slack.InteractionCallback) {
	channel := ""
	settings, err := s.settingsStorage.GetScheduleSettingsForTeam(payload.Team.ID)
	if err != nil || settings == nil {
		fmt.Printf("Print error %+v", err)
	}
	fmt.Printf("SelectedChannel: %+v\n", settings)
	if settings != nil {
		channel = settings.ChannelId
	}
	view, err := handleAppHomeTabChannelPicker(channel)
	_, err = s.client.PublishView(payload.User.ID, view, "")
	if err != nil {
		fmt.Printf("Print error %+v", err)
	}
}

func (s *KudosService) removeKudos(payload slack.InteractionCallback, action *slack.BlockAction) {
	kudosId, err := strconv.Atoi(action.Value)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	err = s.kudosStorage.DeleteKudos(kudosId)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	kudos, _ := s.kudosStorage.GetKudosByUser(payload.User.ID)
	view, err := handleAppHomeTabWithKudosList(kudos)
	_, err = s.client.PublishView(payload.User.ID, view, "")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func (s *KudosService) setSchedule(payload slack.InteractionCallback) {
	err := s.scheduleStorage.SetSchedule(gokudos.Schedule{
		TeamId:      payload.Team.ID,
		ChannelId:   payload.Channel.ID,
		ScheduleId:  "",
		ScheduledAt: time.Now().Add(time.Minute).Unix(),
	})
	if err != nil {
		fmt.Printf("Error while setting a scheduled message: %v\n", err)
	}
}

func (s *KudosService) setChannel(payload slack.InteractionCallback, action *slack.BlockAction) {
	settings := gokudos.Settings{
		TeamId:    payload.Team.ID,
		ChannelId: action.SelectedChannel,
	}
	err := s.settingsStorage.SetScheduleSettings(settings)
	if err != nil {
		fmt.Printf("Print error %+v", err)
	}
	view, err := handleAppHomeTabChannelPicker(settings.ChannelId)
	_, err = s.client.PublishView(payload.User.ID, view, "")
	if err != nil {
		fmt.Printf("Print error %+v", err)
	}
}
