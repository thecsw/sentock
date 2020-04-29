package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// hello is a simple endpoint that returns the message "hello, world"
func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(webError{Msg: "hello, world"})
}

// bobby is a method that deletes all entries in the raw sentiments table of
// the db that are older than 24 hours.  It cleans them up every hour.
func bobby() {
	for {
		log.Infoln("deleting raw sentiments... (anything from before 24 hours ago)")
		err := Elephant.deleteOldRaws(time.Now().UTC().Add(-24 * time.Hour).Unix())
		if err != nil {
			log.Infoln("Problem deleting raw sentiments: " + err.Error())
		}
		// sleep for an hour:
		time.Sleep(time.Hour)
	}
}

// main is the entrypoint. It initializes a postgres (DB) connection
// through Elephant and deploys a web server
func main() {
	log.Infof("Gratiously waiting until other microservices become operational...")
	time.Sleep(5 * time.Second)
	log.Infoln("[DONE]")

	// Connecting to the DB
	log.Infof("Connecting to the DB... ")
	var err error
	db, err = Elephant.connectDB(connectionString)
	if err != nil {
		panic(err)
	}
	log.Infoln("[DONE]")

	// Adding new tables
	log.Infof("Migrating tables... ")
	err = Elephant.autoMigrate()
	if err != nil {
		panic(err)
	}
	log.Infoln("[DONE]")

	// launch goroutine to clean up old raw sentiments
	// GO BOBBY! YOU CAN DO IT
	// I BELIEVE IN YOU BOBBY
	// YOUR FAMILY IS PROUD BOBBY
	// GO GET'EMT TIGER
	// ...
	// and bobby dropped tables
	// relevant:
	go bobby()

	// Create new router
	log.Infof("Creating our HTTP router... ")
	myRouter := mux.NewRouter()
	log.Infoln("[DONE]")

	log.Infof("Adding authorization middleware... ")
	myRouter.Use(authMiddleware)
	log.Infoln("[DONE]")

	log.Infof("Adding HTTP routes... ")
	myRouter.HandleFunc("/", hello).Methods(http.MethodGet)
	myRouter.HandleFunc("/companies", getCompanies).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments", getSentiments).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments", addRawSentiment).Methods(http.MethodPost)
	myRouter.HandleFunc("/sentiments/raw", getRawSentiments).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments/latest", getLatestSentiment).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments/market", getMarketSentiments).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments/averages", addWindowAverages).Methods(http.MethodPost)
	log.Infoln("[DONE]")

	// start listening
	listenAddr := "0.0.0.0:10000"
	log.Infof("Listening on %s... ", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, myRouter))
}
