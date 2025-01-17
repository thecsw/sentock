package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// addRawSentiment is an endpoint that takes a payload of the layout sentimentPayload and adds
// the raw sentiment to the db
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

// getMarketSentiments returns the sentiments since last market open
// until the last close or now, the last one is whichever closer
func getMarketSentiments(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	company := u.Get("company")
	apiType := u.Get("api_type")
	if company == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(webError{Msg: "Received empty params."})
		return
	}
	now := time.Now().UTC()
	// Set the default start and end to today's 14:30 -> 21:00
	start := time.Date(now.Year(), now.Month(), now.Day(), 14, 30, 0, 0, now.Location())
	end := start.Add(6*time.Hour + 30*time.Minute)
	// Handle weekends
	switch {
	// If we are on Saturday or it's before 14:30, grab the day before 14:30 -> 21:00
	case now.Weekday().String() == "Saturday" || now.Before(start):
		start = start.AddDate(0, 0, -1)
		end = end.AddDate(0, 0, -1)
	// If we are on Sunday, remove 2 days
	case now.Weekday().String() == "Sunday":
		start = start.AddDate(0, 0, -2)
		end = end.AddDate(0, 0, -2)
	// If we are before the market close, set it to now
	case now.Hour() < 21:
		end = now
	// Default is do nothing and leave previous values
	default:
		break
	}
	// Get the averaged data:
	averages, err := Elephant.getAverages(company, int(end.Unix()+1), int(start.Unix()-1))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Error getting sentiments: " + err.Error()})
		return
	}
	// PARSE CSV IF REQUESTED AND RETURN
	if apiType == "csv" {
		w.WriteHeader(http.StatusOK)
		csv.NewWriter(w).WriteAll(averagesToString(averages))
		return
	}

	points := make([][]float64, len(averages))
	for ind, data := range averages {
		points[ind] = []float64{float64(data.Unix), data.Average}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(points)
}

// addWindowAverages takes a windowAvePayload and adds the averaged sentiments inside to the db
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

// getRawSentiments retrieves the raw sentiments from a specific company between
// two specific times from the db
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

// getLatestSentiment retrieves the latest unix timestamp from the averaged sentiments in the db
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

// getCompanies is an endpoint that retrieves a list of the unique companies in the db
func getCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := Elephant.getCompanies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(webError{Msg: "Error getting companies: " + err.Error()})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}

// getSentiments is an endpoint to return processed
func getSentiments(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query()
	company := u.Get("company")
	before := u.Get("before")
	after := u.Get("after")
	apiType := u.Get("api_type")
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
	// PARSE CSV IF REQUESTED AND RETURN
	if apiType == "csv" {
		w.WriteHeader(http.StatusOK)
		csv.NewWriter(w).WriteAll(averagesToString(averages))
		return
	}

	points := make([][]float64, len(averages))
	for ind, data := range averages {
		points[ind] = []float64{float64(data.Unix), data.Average}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(points)
}

// averagesToString converts averages to CSV encoded strings
func averagesToString(averages []Average) [][]string {
	csvData := make([][]string, len(averages))
	for ind, data := range averages {
		csvData[ind] = []string{
			strconv.FormatInt(int64(data.Unix), 10),
			strconv.FormatFloat(data.Average, 'f', -1, 64),
		}
	}
	return csvData
}
