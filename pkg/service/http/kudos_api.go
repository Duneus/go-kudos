package http

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/gorilla/mux"
)

type KudosApi struct {
	service gokudos.KudosService
}

func NewKudosApi(service gokudos.KudosService) *KudosApi {
	return &KudosApi{service: service}
}

func (api *KudosApi) Mount(router *mux.Router) {
	router.HandleFunc("/events-endpoint", api.service.Forward)
	router.HandleFunc("/addKudos", api.service.AddKudos)
	router.HandleFunc("/postKudos", api.service.PublishKudos)
	router.HandleFunc("/interactivity", api.service.HandleInteractivity)
}
