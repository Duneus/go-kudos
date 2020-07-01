package main

import (
	"fmt"
	"github.com/Duneus/go-kudos/pkg/config"
	"github.com/Duneus/go-kudos/pkg/inmem"
	"github.com/Duneus/go-kudos/pkg/service"
	"github.com/slack-go/slack"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	api := slack.New(cfg.BotOAuthToken)
	kudosStorage := inmem.NewKudosStorage()
	kudosService := service.NewKudosService(
		kudosStorage,
		cfg,
		api,
	)

	http.HandleFunc("/events-endpoint", kudosService.Handler)

	fmt.Println("[INFO] Server listening")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Server stopped immediately: %v", err)
	}
}
