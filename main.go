package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ukgm/internal/app/config"
	"ukgm/internal/app/config/env"
	firebase_api "ukgm/internal/app/config/firebase"
	"ukgm/internal/app/ukgm/model/request"
	"ukgm/internal/app/ukgm/service"
	webapiv1 "ukgm/internal/app/ukgm/webApiV1"
	login_handler "ukgm/internal/app/ukgm/webApiV1/loginHandler"
	"ukgm/internal/pkg/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func init() {
	config.SetupLogger()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})

	firebaseApi := firebase_api.NewFirebaseApi()

	middlewareSrv := request.NewMiddleware(firebaseApi)

	r.Use(cors.Handler)
	r.Use(middlewareSrv.CommonMiddleware)
	r.Use(middlewareSrv.JWTMiddleware)
	// r.Use(middlewareSrv)

	newService, err := service.ServiceNewDeclare(firebaseApi)
	if err != nil {
		log.Panic(err)
	}
	defer newService.Repo.Client.Disconnect(context.TODO())

	loginApi := login_handler.NewLoginApi(firebaseApi)

	newApi := webapiv1.NewApi(newService, loginApi)
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
