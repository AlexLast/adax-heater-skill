package main

import (
	"net/http"
	"os"
	"time"

	"github.com/alexlast/adax-heater-skill/internal/adax"
	"github.com/alexlast/adax-heater-skill/internal/alexa"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
		},
	})
}

func main() {
	// Build new alexa context
	alexa := new(alexa.Context)

	// Define the Adax client with an HTTP client
	alexa.Adax = &adax.Client{
		Config: &adax.Config{},
		HTTP: &http.Client{
			Timeout: time.Second * 3,
		},
	}

	// Unmarshal configuration into Adax.Config type
	err := envconfig.Process("adax", alexa.Adax.Config)

	if err != nil {
		log.Fatalln(err)
	}

	// Start the lambda
	lambda.Start(alexa.Handler)
}
