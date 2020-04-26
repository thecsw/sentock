package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	serverIP            = "http://server:10000"
	globpostAves        = serverIP + "/sentiments/averages"
	globlatestsentiment = serverIP + "/sentiments/latest"
	globrawsentiments   = serverIP + "/sentiments/raw"

	wg sync.WaitGroup
)

// windowSmoothing is a function that takes a list of datapoints (unix timestamp, sentiment value),
// a starting value to calculate, and a window size and
// returns a list of datapoints (unix timestamp, ave sentiment for past hour)
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

func handleRealtime(company, postAves, latestsentiment, rawsentiments string) {
	defer wg.Done()
	buffer := make([][]float64, 0)
	tries := 0
	for {
		log.Infof("[%s]: DOING A SWEEP", company)
		waketime := time.Now()
		//check how far calculations have come (assumes everything up to the most recent empty minute has already been calculated correctly):
		resp, err := http.Get(latestsentiment + "?company=" + company)
		//ensure we got something from the server
		if err != nil || resp == nil || resp.StatusCode != 200 {
			tries++
			log.Error("failed to get latest sentiment for company: ", company, " on try #", tries, ", trying again in a bit...")
			log.Error("[ERROR]:", err)
			if resp != nil {
				log.Error("[STATUS]:", resp.StatusCode)
			}
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}
		//decode result:
		var result int64
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&result)
		if err != nil {
			log.Error("Failed to decode latest sentiment query response")
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}
		//get list of new sentiments:
		values := url.Values{}
		values.Set("company", company)
		values.Set("before", strconv.FormatInt(waketime.Unix(), 10))
		values.Set("after", strconv.FormatInt(result-3600, 10)) //get the data from an hour before last entry (just in case buffer is empty)
		resp, err = http.Get(rawsentiments + "?" + values.Encode())
		// Don't forget to check resp.StatusCode for 200
		if err != nil || resp.StatusCode != 200 {
			log.Error("Couldn't access rawsentiments from: " + rawsentiments + "?" + values.Encode())
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}
		var rawsents [][]float64
		decoder = json.NewDecoder(resp.Body)
		err = decoder.Decode(&rawsents)
		//ensure rawsents has no repeats and each element has the expected size of 3 (Unix, Sentiment, TweetID)
		//fmt.Println("Length of rawSentiments returned: ", len(rawsents))
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
		//calculate and fill in averages needed: (uses result (the last unix timestamp calculated) as the front
		averages := windowSmoothing(buffer, result, 3600)
		if len(averages) == 0 {
			time.Sleep(2 * time.Minute)
			continue
		}
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
			log.Error("Couldn't jsonify the averages list!")
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil || resp == nil || resp.StatusCode != 200 {
			log.Error("couldn't push average sentiments: [ERROR]:", err)
			if resp != nil {
				log.Error("[STATUS]:", resp.StatusCode, " body: ", resp.Body)
			}
			resp.Body.Close()
			tries++
			time.Sleep(time.Duration(tries) * time.Second)
			continue
		}
		resp.Body.Close()
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
	}
}

func main() {
	log.Infoln("Gratiously waiting until other microservices become operational...")
	// From https://play.golang.org/
	// Clear the screen by printing \x0c.
	const col = 30
	bar := fmt.Sprintf("\x0c[%%-%vs]", col)
	for i := 0; i < col; i++ {
		log.Printf(bar, strings.Repeat("=", i)+">")
		time.Sleep(100 * time.Millisecond)
	}
	log.Printf(bar+" Done!", strings.Repeat("=", col))
	//do window averaging calculations:
	companies := []string{
		"McDonalds",
		"Fedex",
		"Chipotle",
		"Microsoft",
		"Disney",
		"Tesla",
		"Google",
		"Facebook",
		"Amazon",
	}
	for ind := 0; ind < len(companies); ind++ {
		wg.Add(1)
		log.Infof("Starting goroutine for %s", companies[ind])
		go handleRealtime(companies[ind], globpostAves, globlatestsentiment, globrawsentiments)
	}
	wg.Wait()
	os.Exit(1)
}
