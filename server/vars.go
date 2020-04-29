package main

var (
	cachedNames map[string]uint = make(map[string]uint)
)

// webError is an interface for HTTP error return
// codes
type webError struct {
	Msg string `json:"msg"`
}

// sentimentPayload is what server expects to receive
// in a request body when adding a new sentiment value
type sentimentPayload struct {
	Company   string  `json:"company"`
	TweetID   string  `json:"tweet_id"`
	Sentiment float64 `json:"sentiment"`
	Unix      int     `json:"unix"`
}

// windowAvePayload is the body request we want
// whan adding now average values (usually josh)
type windowAvePayload struct {
	Company  string    `json:"company"`
	Unix     []int     `json:"unix"`
	Averages []float64 `json:"average"`
}
