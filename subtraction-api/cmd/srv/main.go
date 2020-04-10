package main

import (
	"log"
	"github.com/rs/cors"
	"net/http"
	"strconv"
	"math-api/base-api"
	"math-api/subtraction-api/pkg/operate"
)

var config *base_api.Config

func init() () {
	var err error

	config, err = base_api.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() () {
	subtractionOperation := &base_api.Operation{
		Operate:                operate.Subtract,
		ExpectedArgumentLength: 2,
	}

	a := base_api.New(config)

	r := base_api.CreateRouter(a, subtractionOperation)

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
