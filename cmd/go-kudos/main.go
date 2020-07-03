package main

import (
	"fmt"
	"github.com/Duneus/go-kudos/pkg/config"
	"github.com/Duneus/go-kudos/pkg/service"
	http2 "github.com/Duneus/go-kudos/pkg/service/http"
	"github.com/Duneus/go-kudos/pkg/sqlite"
	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	api := slack.New(cfg.BotOAuthToken)
	db, err := sqlite.NewGorm(cfg.SqliteFilePath)
	if err != nil {
		panic(err)
	}
	kudosPersistentStorage := sqlite.NewKudosStorage(db)
	kudosService := service.NewKudosService(
		kudosPersistentStorage,
		cfg,
		api,
	)

	kudosApi := http2.NewKudosApi(kudosService)

	router := mux.NewRouter()
	kudosApi.Mount(router)


	fmt.Println("[INFO] Server listening")

	if err := http.ListenAndServe(":3000", router); err != nil {
		fmt.Printf("Server stopped immediately: %v", err)
	}
}
