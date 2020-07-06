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
)

var _ gokudos.KudosService = &KudosService{}

type KudosService struct {
	kudosStorage gokudos.KudosStorage
	cfg          config.Config
	client       *slack.Client
}

func NewKudosService(
	storage gokudos.KudosStorage,
	cfg config.Config,
	client *slack.Client,
) *KudosService {
	return &KudosService{
		kudosStorage: storage,
		cfg:          cfg,
		client:       client,
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
		message = message + k.Message + "\n"
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
		switch action.ActionID {
		case "show_kudos":
			kudos, _ := s.kudosStorage.GetKudosByUser(payload.User.ID)
			view, err = handleAppHomeTabWithKudosList(kudos)
		case "show_all_kudos":
			kudos, _ := s.kudosStorage.GetAllKudosInTeam(payload.Team.ID)
			view, err = handleAppHomeTabWithKudosList(kudos)
		case "remove_kudos":
			s.kudosStorage.DeleteKudos(action.Value)
			kudos, _ := s.kudosStorage.GetKudosByUser(payload.User.ID)
			view, err = handleAppHomeTabWithKudosList(kudos)
		case "hide_kudos":
			view, err = handleAppHomeTab()
		default:
			view, err = handleAppHomeTab()
		}
		_, _ = s.client.PublishView(payload.User.ID, view, "")
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
