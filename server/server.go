package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	befr, err := strconv.Atoi(before)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Error converting before to int: " + err.Error()})
		return
	}
	aftr, err := strconv.Atoi(after)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Error converting after to int: " + err.Error()})
		return
	}
	//set moving average window size:
	window := 3600 //window of 1 hour
	// call Elephant.getSentiments(company, before, after)
	sentiments, err := Elephant.getSentiments(company, befr, aftr-window)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Error getting sentiments: " + err.Error()})
		return
	}
	if len(sentiments) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([][]float64{[]float64{}})
		return
	}
	// Do the window smoothing algorithm
	front := (aftr / 60) * 60 //front starts at the minute of the furthest back value wanted
	points := windowSmoothing(sentiments, front, window)
	// Submit the new one
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(points) // [[...],[...],[...], ...]
}

func windowSmoothing(sentiments []Sentiment, front, window int) [][]float64 {
	var end = make([][]float64, 0, 100)
	back := front - window
	var backind, frontind int
	sum := 0.0
	size := 0
	for front < sentiments[len(sentiments)-1].Unix {
		//find new backind (subtracting as we go)
		for sentiments[backind].Unix < back {
			sum -= sentiments[backind].Sentiment
			backind++
			size--
		}
		//find new frontind (subtracting as we go)
		for sentiments[frontind].Unix < front {
			sum += sentiments[frontind].Sentiment
			frontind++
			size++
		}
		//get new ave, assign to end
		end = append(end, []float64{float64(front), sum / float64(size)})
		//increment front (and back)
		front += 60
		back += 60
	}
	return end
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
	myRouter.HandleFunc("/sentiments", getSentiments).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments", addSentiment).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
