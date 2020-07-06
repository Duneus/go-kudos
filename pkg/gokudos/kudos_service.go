package gokudos

import (
	"net/http"
)

type KudosService interface {
	Forward(rw http.ResponseWriter, r *http.Request)
	AddKudos(rw http.ResponseWriter, r *http.Request)
	PublishKudos(rw http.ResponseWriter, r *http.Request)
	HandleInteractivity(rw http.ResponseWriter, r *http.Request)
}
