package gokudos

import (
	"net/http"
)

type KudosService interface {
	Handler(rw http.ResponseWriter, r *http.Request)

	SendMessage(channelId string, message string) error
	ForwardKudos()
	PublishKudos()
}
