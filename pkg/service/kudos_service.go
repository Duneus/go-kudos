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

type ActionType string

const (
	ShowKudosView         ActionType = "show_kudos_view"
	ShowAllKudosView      ActionType = "show_all_kudos_view"
	ShowSchedulingView    ActionType = "show_scheduling_view"
	ShowChannelSelectView ActionType = "show_channel_select_view"

	SetSchedule ActionType = "set_schedule"
	SetChannel  ActionType = "set_channel"
	RemoveKudos ActionType = "remove_kudos"
)

type InteractivityHandler func(payload slack.InteractionCallback, action *slack.BlockAction) error

type HandlerMap map[ActionType]InteractivityHandler

type KudosService struct {
	kudosStorage    gokudos.KudosStorage
	scheduleStorage gokudos.ScheduleStorage
	settingsStorage gokudos.SettingsStorage
	cfg             config.Config
	client          *slack.Client
	handlerMap      HandlerMap
}

func NewKudosService(
	kudosStorage gokudos.KudosStorage,
	scheduleStorage gokudos.ScheduleStorage,
	settingsStorage gokudos.SettingsStorage,
	cfg config.Config,
	client *slack.Client,
	opts ...func(HandlerMap),
) *KudosService {
	handlerMap := make(HandlerMap)

	for _, opt := range opts {
		opt(handlerMap)
	}
	return &KudosService{
		kudosStorage:    kudosStorage,
		scheduleStorage: scheduleStorage,
		settingsStorage: settingsStorage,
		cfg:             cfg,
		client:          client,
		handlerMap:      handlerMap,
	}
}

func WithActionHandler(actionType ActionType, handler InteractivityHandler) func(handlerMap HandlerMap) {
	return func(handlerMap HandlerMap) {
		handlerMap[actionType] = handler
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
		err := s.handlerMap[ActionType(action.ActionID)](payload, action)
		if err != nil {
			fmt.Printf("No handler found received action: %v", err)
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

func NewShowMyKudosViewHandler(storage gokudos.KudosStorage, client *slack.Client) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		kudos, err := storage.GetKudosByUser(payload.User.ID)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		view, err := handleAppHomeTabWithKudosList(kudos)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		_, err = client.PublishView(payload.User.ID, view, "")
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		return nil
	}
}

func NewShowAllKudosViewHandler(storage gokudos.KudosStorage, client *slack.Client) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		kudos, err := storage.GetAllKudosInTeam(payload.Team.ID)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		view, err := handleAppHomeTabWithKudosList(kudos)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		_, err = client.PublishView(payload.User.ID, view, "")
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		return nil
	}
}

func NewShowSchedulingViewHandler(client *slack.Client) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		view, err := handleAppHomeSchedulingTab()
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		_, err = client.PublishView(payload.User.ID, view, "")
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		return nil
	}
}

func NewShowChannelSelectViewHandler(storage gokudos.SettingsStorage, client *slack.Client) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		channel := ""
		settings, err := storage.GetScheduleSettingsForTeam(payload.Team.ID)
		if err != nil || settings == nil {
			fmt.Printf("Print error %+v", err)
		}
		fmt.Printf("SelectedChannel: %+v\n", settings)
		if settings != nil {
			channel = settings.ChannelId
		}
		view, err := handleAppHomeTabChannelPicker(channel)
		_, err = client.PublishView(payload.User.ID, view, "")
		if err != nil {
			fmt.Printf("Print error %+v", err)
		}

		return nil
	}
}

func NewRemoveKudosHandler(storage gokudos.KudosStorage, client *slack.Client) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		kudosId, err := strconv.Atoi(action.Value)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		err = storage.DeleteKudos(kudosId)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		kudos, _ := storage.GetKudosByUser(payload.User.ID)
		view, err := handleAppHomeTabWithKudosList(kudos)
		_, err = client.PublishView(payload.User.ID, view, "")
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		return nil
	}
}

func NewSetScheduleHandler(storage gokudos.ScheduleStorage) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		err := storage.SetSchedule(gokudos.Schedule{
			TeamId:      payload.Team.ID,
			ChannelId:   payload.Channel.ID,
			ScheduleId:  "",
			ScheduledAt: time.Now().Add(time.Minute).Unix(),
		})
		if err != nil {
			fmt.Printf("Error while setting a scheduled message: %v\n", err)
		}

		return nil
	}
}

func NewSetChannelHandler(storage gokudos.SettingsStorage, client *slack.Client) InteractivityHandler {
	return func(payload slack.InteractionCallback, action *slack.BlockAction) error {
		settings := gokudos.Settings{
			TeamId:    payload.Team.ID,
			ChannelId: action.SelectedChannel,
		}
		err := storage.SetScheduleSettings(settings)
		if err != nil {
			fmt.Printf("Print error %+v", err)
		}
		view, err := handleAppHomeTabChannelPicker(settings.ChannelId)
		_, err = client.PublishView(payload.User.ID, view, "")
		if err != nil {
			fmt.Printf("Print error %+v", err)
		}
		return nil
	}
}
