package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"net/http"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New("xoxb-1195463250596-1189541187811-xZTcNBCw9RS5hmI6hiJpPSJX")

func main() {
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: "hQbvWoUqPoQE24brikPBPCC8"}))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			fmt.Printf("Event inner type: %s\n", innerEvent.Data)
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			case *slackevents.MessageEvent:
				err := HandleMessageEvents(ev)
				if err != nil {
					fmt.Printf("Error handling message: %v", err)
				}
			case *slackevents.MessageAction:
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Server stopped immediately: %v", err)
	}
}

func HandleMessageEvents(event *slackevents.MessageEvent) error {
	if event.BotID == "" {
		api.PostMessage(event.Channel, slack.MsgOptionText("Hello, thanks", false))
		fmt.Printf("message sender: %s\n", event.User)
		fmt.Printf("message sender id: %s\n", event.Username)
		fmt.Printf("bot: %s\n", event.BotID)
	}

	return nil
}