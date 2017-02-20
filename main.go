package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/todoapi/models"
	_ "github.com/todoapi/models/mongo"
	_ "github.com/todoapi/models/sqlite"
)

var (
	dbModelName = flag.String("dbmodel", "mongo", "Database Backend Model")
	debug       = flag.Bool("debug", false, "Debug")
	model       models.Model
)

func getEnv() {
	logrus.Info("Starting The TODO API")
	flag.Parse()
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var err error
	model, err = models.GetDBModel(*dbModelName)
	if err != nil {
		logrus.Fatalf("Failed to connect to database '%s': %v", *dbModelName, err)
	}
}

func main() {
	getEnv()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
