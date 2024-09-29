package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/api"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/cfg"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/database"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/repository"
	"github.com/etloreq/yandex-practicum-arch-tasks/sprint3/microservice-telemetry/internal/service"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := database.MustDB(config.Database)
	defer db.Close()

	repo := repository.New(db, config.Org, config.Bucket)
	srv := service.New(repo)

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer(srv)

	r := chi.NewMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(api.NewStrictHandler(server, nil), r)

	s := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
	}

	log.Println("start server")

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
