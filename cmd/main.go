package main

import (
	"api/internal/application"
	"api/internal/infrastructure"
	"api/internal/transport/http"
	"api/internal/transport/websocket"

	"api/pkg/configuration"
	"api/pkg/postgresql"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgresql.Connect(&cfg.Postgresql)
	if err != nil {
		log.Fatal(err)
	}

	repository := infrastructure.New(db)
	useCase := application.New(repository)

	ws := websocket.New(useCase.Users, useCase.Chats)
	httpServer := http.New(useCase, ws)

	//
	go ws.Run()
	httpServer.HandleWS("/ws", ws.GetServer())
	//

	go func() {
		if err := httpServer.Run(cfg.HttpSocket); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	_, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
}
