package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ukgm/internal/app/config/env"
	"ukgm/internal/app/ukgm/model/request"
	"ukgm/internal/app/ukgm/service"
	webapiv1 "ukgm/internal/app/ukgm/webApiV1"
	simple "ukgm/internal/app/ukgm/webApiV1/simpleHandler"
	"ukgm/internal/pkg/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})

	r.Use(cors.Handler)
	r.Use(request.CommonMiddleware)

	newService, err := service.ServiceNewDeclare()
	if err != nil {
		log.Panic(err)
	}

	newSimpleApi := simple.NewSimpleApi(newService)

	newApi := webapiv1.NewApi(newService, newSimpleApi)
	r.Route("/api/v1/", newApi.Init)

	appHost := env.Get("PORT")
	log.Println("Running at : ", appHost)
	appHost = ":" + appHost

	/*
		Gracefully stop application
		reference : https://medium.com/@kpbird/golang-gracefully-stop-application-23c2390bb212
	*/

	s := server.New(&http.Server{
		Addr:    appHost,
		Handler: r,
	})
	s.Start()
	defer s.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Server is shutting down")
}
