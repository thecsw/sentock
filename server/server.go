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
		//some sort of return error
		return
	}
	aftr, err := strconv.Atoi(after)
	if err != nil {
		//some sort of return error
		return
	}
	//set moving average window size:
	window := 3600 //window of 1 hour
	// call Elephant.getSentiments(company, before, after)
	sentiments, err := Elephant.getSentiments(company, befr, aftr-window)
	if err != nil {
		//some sort of error handling
		return
	}
	// Do the window smoothing algorithm
	front := (aftr / 60) * 60 //front starts at the minute of the furthest back value
	points := windowSmoothing(sentiments, front, window)
	// Submit the new one
}

func windowSmoothing(sentiments []Sentiment, front, window int) []dataPoint {
	return nil
}

/*
	window = 60*60 #window size of an hour

    c = conn.cursor()
    c.execute("""
    SELECT timestamp, sentiment FROM companies WHERE company = %s AND timestamp < %s AND timestamp > %s ORDER BY timestamp DESC;
    """, (company, int(before), (int(after)-window),))
    result = c.fetchall()
    c.close()
    if (len(result) < 1):
        return []
    #moving average (of time period of size [window] seconds):
    latest = (result[0][0]//60)*60 #anchor on the minute mark
    end_result = []
    start = 0
    while latest > int(after):
        (window_ave, start) = get_window_ave(result, window, latest, start)
        if window_ave is not None:
            end_result.append((latest,window_ave))
        latest-=60
    return end_result

def get_window_ave(vals, window, latest, start):
    (values, startNew) = get_vals_of_interval(vals, latest, latest-window, start)
    if (len(values)<=5):
        return (None, startNew)
    sumVal = 0
    for val in values:
        sumVal += val[1]
    return (sumVal / len(values), startNew)

def get_vals_of_interval(vals, before, after, start):
    ans = []
    startNew = start
    for val in vals[start:]:
        if (val[0] > before - 60):
            startNew += 1
        if (val[0] < before and val[0] > after):
            ans.append(val)
        if (val[0] < after):
            break
    return (ans, startNew)
*/

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
