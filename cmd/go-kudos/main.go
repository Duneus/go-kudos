package main

import (
	"fmt"
	"github.com/Duneus/go-kudos/pkg/config"
	"github.com/Duneus/go-kudos/pkg/service"
	http2 "github.com/Duneus/go-kudos/pkg/service/http"
	slack2 "github.com/Duneus/go-kudos/pkg/slack"
	"github.com/Duneus/go-kudos/pkg/sqlite"
	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	api := slack.New(cfg.BotOAuthToken)
	db, err := sqlite.NewGorm(cfg.SqliteFilePath)
	//sqlite.Migrate(db)
	if err != nil {
		panic(err)
	}
	kudosPersistentStorage := sqlite.NewKudosStorage(db)
	scheduleStorage := slack2.NewScheduleStorage(api)
	settingsStorage := sqlite.NewSettingsStorage(db)
	kudosService := service.NewKudosService(
		kudosPersistentStorage,
		scheduleStorage,
		settingsStorage,
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
