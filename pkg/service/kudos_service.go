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

func (s *KudosService) Handler(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := s.parseMessage(body)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.handleEvent(eventsAPIEvent, body, w)
}

func (s *KudosService) SendMessage(channelId string, message string) error {
	_, _, err := s.client.PostMessage(channelId, slack.MsgOptionText(message, false))

	return err
}


func (s *KudosService) ForwardKudos() {
	panic("implement me")
}

func (s *KudosService) PublishKudos() {
	panic("implement me")
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
		fmt.Printf("Event inner type: %s\n", innerEvent.Data)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			s.SendMessage(ev.Channel, "Yes, hello.")
		case *slackevents.MessageEvent:
			err := s.handleMessageEvents(ev)
			if err != nil {
				fmt.Printf("Error handling message: %v", err)
			}
		case *slackevents.MessageAction:
		}
	}
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
		s.SendMessage(event.Channel, "Hello, thanks")
		fmt.Printf("message sender: %s\n", event.User)
		fmt.Printf("message sender id: %s\n", event.Username)
		fmt.Printf("bot: %s\n", event.BotID)
	}

	return nil
}
