package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type webError struct {
	Msg string `json:"msg"`
}

type sentimentPayload struct {
	Company   string  `json:"company"`
	TweetID   string  `json:"tweet_id"`
	Sentiment float64 `json:"sentiment"`
	Unix      int     `json:"unix"`
}

func addSentiment(w http.ResponseWriter, r *http.Request) {
	payload := &sentimentPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: err.Error()})
		return
	}
	res, err := Elephant.createSentiment(
		payload.TweetID,
		payload.Unix,
		payload.Sentiment,
		payload.Company,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*res)
}

func getSentiments(w http.ResponseWriter, r *http.Request) {
	company := mux.Vars(r)["company"]
	before := mux.Vars(r)["before"]
	after := mux.Vars(r)["after"]
	// make before and after ints
	// call Elephant.getSentiments(company, before, after)
	// Do the window smoothing algorithm
	// Submit the new one
}

func main() {
	// Connecting to the DB
	fmt.Printf("Connecting to the DB... ")
	var err error
	db, err = Elephant.connectDB(connectionString)
	if err != nil {
		panic(err)
	}
	fmt.Println("[DONE]")

	// Adding new tables
	fmt.Printf("Migrating tables... ")
	err = Elephant.autoMigrate()
	if err != nil {
		panic(err)
	}
	fmt.Println("[DONE]")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/sentiments", addSentiment).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
