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

type windowAvePayload struct {
	Company  string    `json:"company"`
	Unix     []int     `json:"unix"`
	Averages []float64 `json:"average"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(webError{Msg: "hello, world"})
}

func addRawSentiment(w http.ResponseWriter, r *http.Request) {
	// get sentiment data:
	payload := &sentimentPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Failed decoding: " + err.Error()})
		return
	}
	// ensure reasonable inputs:
	if payload.TweetID == "" || payload.Unix == 0 || payload.Company == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(webError{Msg: "Received empty params."})
		return
	}
	// add sentiment to db:
	res, err := Elephant.createSentiment(
		payload.TweetID,
		payload.Unix,
		payload.Sentiment,
		payload.Company,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Failed committing to DB: " + err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*res)
}

func addWindowAverages(w http.ResponseWriter, r *http.Request) {
	// get sentiment data:
	payload := &windowAvePayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Failed decoding: " + err.Error()})
		return
	}
	// ensure reasonable inputs:
	if payload.Company == "" || len(payload.Unix) != len(payload.Averages) || len(payload.Unix) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(webError{Msg: "Received empty params."})
		return
	}
	// add average sentiments to db:
	_, err = Elephant.createWindowAverage(
		payload.Unix,
		payload.Averages,
		payload.Company,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Failed committing to DB: " + err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(*res)
}

func getRawSentiments(w http.ResponseWriter, r *http.Request) {
	// get inputs:
	u := r.URL.Query()
	company, before, after := u.Get("company"), u.Get("before"), u.Get("after")
	// ensure valid inputs:
	if before == "" || after == "" || company == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(webError{Msg: "Received empty params."})
		return
	}
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
	// call Elephant.getSentiments(company, before, after)
	sentiments, err := Elephant.getSentiments(company, befr, aftr)
	// make and fill in array containing just the unix timestamp and the sentiment value:
	points := make([][]float64, len(sentiments))
	for ind, sent := range sentiments {
		temp, _ := strconv.ParseFloat(sent.TweetID, 64)
		//points[ind] = []float64{sent.Unix, sent.Sentiment, float64(strconv.ParseInt(sent.TweetID,10,0))}
		points[ind] = []float64{float64(sent.Unix), sent.Sentiment, temp}
	}
	// return the list of datapoints:
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(points)
}

func getLatestSentiment(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	company := u.Get("company")
	if company == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(webError{Msg: "Received empty company."})
		return
	}
	answer, err := Elephant.getLatestAverageSentiment(company)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting last average sentiment: " + err.Error())
		json.NewEncoder(w).Encode(webError{Msg: "Error getting last average sentiment: " + err.Error()})
		return
	}
	// return the list of datapoints:
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answer)
}

// getSentiments is a blablabla
func getSentiments(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	company := u.Get("company")
	before := u.Get("before")
	after := u.Get("after")
	fmt.Println("|C:", company, "|B:", before, "|A:", after)
	if before == "" || after == "" || company == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(webError{Msg: "Received empty params."})
		return
	}
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
	//get the averaged data:
	averages, err := Elephant.getAverages(company, befr, aftr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Error getting sentiments: " + err.Error()})
		return
	}
	points := make([][]float64, len(averages))
	for ind, data := range averages {
		points[ind] = []float64{float64(data.Unix), data.Average}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(points)
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

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/", hello).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments", getSentiments).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments", addRawSentiment).Methods(http.MethodPost)
	myRouter.HandleFunc("/sentiments/raw", getRawSentiments).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments/latest", getLatestSentiment).Methods(http.MethodGet)
	myRouter.HandleFunc("/sentiments/averages", addWindowAverages).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe("0.0.0.0:10000", myRouter))
}
