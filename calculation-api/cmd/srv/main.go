package main

import (
	"log"
	"github.com/rs/cors"
	"net/http"
	"strconv"
	"math-api/calculation-api/pkg/api"
	"math-api/calculation-api/pkg/router"
	"math-api/calculation-api/pkg/configurations"
)

var config *configurations.Config

func init() () {
	var err error

	config, err = configurations.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() () {
	a := api.New(config)

	r := router.CreateRouter(a)

	c := cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowCredentials: config.Cors.AllowCredentials,
	}

	log.Printf("Server listening on port %v.\n", config.Port)
	log.Fatalln(
		http.ListenAndServe(
			":"+strconv.Itoa(config.Port),
			cors.New(c).Handler(r),
		),
	)
}
