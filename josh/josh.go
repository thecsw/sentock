package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	serverIP            = "http://server:10000"
	companiesVals       = serverIP + "/companies"
	globMarketVals      = serverIP + "/sentiments/market"
	globpostAves        = serverIP + "/sentiments/averages"
	globlatestsentiment = serverIP + "/sentiments/latest"
	globrawsentiments   = serverIP + "/sentiments/raw"
)

// windowSmoothing is a function that takes a list of datapoints (unix timestamp, sentiment value),
// a starting value to calculate, and a window size and
// returns a list of datapoints with values: (unix timestamp, ave sentiment for past hour)
func windowSmoothing(raws [][]float64, front, window int64) [][]float64 {
	if len(raws) == 0 {
		log.Error("no data to window smooth")
		return [][]float64{}
	}
	var end = make([][]float64, 0, 100)
	back := front - window
	frontind := len(raws) - 1
	backind := len(raws) - 1
	sum := 0.0
	size := 0
	for front < int64(raws[0][0]) && frontind > -1 {
		//find new frontind (subtracting as we go)
		for int64(raws[frontind][0]) < front && frontind > -1 {
			sum += raws[frontind][1]
			frontind--
			size++
		}
		//find new backind (subtracting as we go)
		for int64(raws[backind][0]) < back && backind > -1 {
			sum -= raws[backind][1]
			backind--
			size--
		}
		//get new average, assign to end
		if size >= 15 {
			end = append([][]float64{[]float64{float64(front), sum / float64(size)}}, end...) //prepend to ensure descending chronological order
		}
		//increment front (and back)
		front += 60
		back += 60
	}
	return end
}

// getLatest queries the server in order to get the latest average sentiment calculated from ip=latestsentiment and company=company
func getLatest(company, latestsentiment string) (int64, error) {
	//check how far calculations have come (assumes everything up to the most recent empty minute has already been calculated correctly):
	resp, err := http.Get(latestsentiment + "?company=" + company)
	//ensure good response:
	if err != nil {
		return -1, errors.New("Error getting latest sentiment from server: " + err.Error())
	}
	if resp == nil {
		return -1, errors.New("Unkown error occured getting latest sentiment from server (nil response)")
	}
	if resp.StatusCode != 200 {
		return -1, errors.New("Got unexpected status code while trying to retrieve latest sentiment: " + strconv.Itoa(resp.StatusCode))
	}
	//decode result:
	var result int64
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return -1, errors.New("Failed to decode latest sentiment query response: " + err.Error())
	}
	return result, nil
}

// getrawSentiments queries the server for the most recent raw sentiments
func getrawSentiments(company, rawSentsIp string, waketime time.Time, latest, windowSize int64) ([][]float64, error) {
	//get list of new sentiments:
	values := url.Values{}
	values.Set("company", company)
	values.Set("before", strconv.FormatInt(waketime.Unix(), 10))
	values.Set("after", strconv.FormatInt(latest-windowSize, 10)) //get the data from an hour before last entry (just in case buffer is empty)
	resp, err := http.Get(rawSentsIp + "?" + values.Encode())
	//ensure good response:
	if err != nil {
		return [][]float64{}, errors.New("Error occured trying to get rawSentiments: " + err.Error())
	}
	if resp == nil {
		return [][]float64{}, errors.New("Unexpected error occured trying to get rawSentiments (nil response)")
	}
	if resp.StatusCode != 200 {
		return [][]float64{}, errors.New("Got unexpected status code while trying to get rawSentiments: " + strconv.Itoa(resp.StatusCode))
	}
	//decode result:
	var rawsents [][]float64
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&rawsents)
	if err != nil {
		return [][]float64{}, errors.New("Error occured while trying to parse rawSentiments response: " + err.Error())
	}
	return rawsents, nil
}

// postAverages posts the processed sentiments to the server so they can be added to the db.
func postAverages(company, postAves string, averages [][]float64) error {
	type windowAvePayload struct {
		Company  string    `json:"company"`
		Unix     []int     `json:"unix"`
		Averages []float64 `json:"average"`
	}
	data := &windowAvePayload{Company: company, Unix: []int{}, Averages: []float64{}}
	for _, val := range averages {
		data.Unix = append(data.Unix, int(val[0]))
		data.Averages = append(data.Averages, val[1])
	}
	jsonAves, err := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, postAves, bytes.NewBuffer(jsonAves))
	req.Header.Add("webkey", os.Getenv("WEB_KEY"))
	if err != nil {
		return errors.New("Couldn't jsonify the averages list when trying to post averages")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return errors.New("Couldn't post the average sentiments: (either nil response or non-nil error: " + err.Error() + ")")
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return errors.New("recieved unexpected status code while trying to post averages: " + strconv.Itoa(resp.StatusCode))
	}
	resp.Body.Close()
	return nil
}

// handleRealtime handles the reatlime analysis of the raw sentiment data.  It is thread safe for different companies.
// It takes in the company and the list of different endpoints to retreive or post data to.
func handleRealtime(company, postAves, latestsentiment, rawsentiments string) {
	buffer := make([][]float64, 0)
	tries := 0
	for {
		log.Infof("[%s]: DOING A SWEEP", company)
		waketime := time.Now()

		//get latest sentiment calculated:
		result, err := getLatest(company, latestsentiment)
		if err != nil {
			log.Error("Company: " + company + ", Try #" + strconv.Itoa(tries) + err.Error())
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}

		//get new unprocessed sentiments:
		rawsents, err := getrawSentiments(company, rawsentiments, waketime, result, 3600)
		if err != nil {
			log.Error("Company: " + company + ", Try #" + strconv.Itoa(tries) + err.Error())
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}

		//if we already have things in the buffer, ensure rawsents has no repeats and each element has the expected size of 3 (Unix, Sentiment, TweetID)
		if len(buffer) > 0 {
			for ind, point := range rawsents {
				if len(point) < 3 || point[2] == buffer[0][2] {
					rawsents = rawsents[0:ind] //since rawsents is garaunteed to be in descending order
					break
				}
			}
		}
		//prepend results to buffer: (so that buffer always remains in descending order)
		buffer = append(rawsents, buffer...)

		//calculate and fill in averages needed: (uses result+60 (the next unix timestamp to calculate) as the front)
		averages := windowSmoothing(buffer, result+60, 3600)
		//if we got no tweets back, proceed like normal
		if len(averages) == 0 {
			time.Sleep(2 * time.Minute)
			continue
		}

		//post resulting calculations to server:
		err = postAverages(company, postAves, averages)
		if err != nil {
			log.Error("Company: " + company + ", Try #" + strconv.Itoa(tries) + err.Error())
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}

		//clean up end of buffer if need be:
		for ind := 0; ind < len(buffer); ind++ {
			if averages[0][0]-buffer[ind][0] > 3600 {
				buffer = buffer[0:ind]
				break
			}
		}
		log.Infof("[%s]: FINISHED MY SWEEP", company)

		//sleep for 2 minutes:
		time.Sleep(2 * time.Minute)

		//reset tries to 3:
		if tries > 3 {
			tries = 3
		}
	}
}

// main queries the server for a list of companies and then starts up the realtime analysis of each company's raw sentiments.
func main() {
	log.Infoln("Gratiously waiting until other microservices become operational...")
	// From https://play.golang.org/
	// Clear the screen by printing \x0c.
	const col = 30
	bar := fmt.Sprintf("\x0c[%%-%vs]", col)
	for i := 0; i < col; i++ {
		log.Printf(bar, strings.Repeat("=", i)+">")
		time.Sleep(1 * time.Second)
	}
	log.Printf(bar+" Done!", strings.Repeat("=", col))

	// Get all companies in the system
	resp, err := http.Get(companiesVals)
	if err != nil {
		panic("Failed to receive initial companies: " + err.Error())
	}
	companies := make([]string, 0, 16)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&companies)
	if err != nil {
		panic("Failed to decode initial companies: " + err.Error())
	}
	log.Println("Received the following companies:", companies)
	for _, company := range companies {
		log.Infof("Starting goroutine for %s", company)
		go handleRealtime(company, globpostAves, globlatestsentiment, globrawsentiments)
	}
	// The stop signal channel
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	log.Println("shutting down")
	os.Exit(0)
}
