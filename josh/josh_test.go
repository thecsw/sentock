package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type joshTestSuite struct {
	suite.Suite
}

func (suite *joshTestSuite) TestAWindowSmoothing() {
	rawdata := make([][]float64, 0, 100)
	// less than 15 values in a window ensures no average is generated:
	for i := 1; i <= 14; i++ {
		rawdata = append([][]float64{[]float64{20 + float64(i), .5}}, rawdata...) // prepending to ensure desc order
	}
	// first window has average of .75
	for i := 1; i <= 20; i++ {
		rawdata = append([][]float64{[]float64{3600 + float64(i), .75}}, rawdata...)
	}
	// second window has average of 0
	for i := 1; i <= 20; i++ {
		rawdata = append([][]float64{[]float64{3660 + float64(i), -.75}}, rawdata...)
	}
	// final value greater than last value considered to help ensure that previous values
	// won't need to be calculated again
	rawdata = append([][]float64{[]float64{3721, 20000}}, rawdata...)
	output := windowSmoothing(rawdata, 3600, 3600)
	suite.Assert().Equal(2, len(output))
	suite.Assert().Equal(2, len(output[0]))
	suite.Assert().Equal(2, len(output[1]))
	suite.Assert().Equal(float64(3660), output[1][0]) // first timestamp correct
	suite.Assert().Equal(float64(3720), output[0][0]) // second timestamp correct
	suite.Assert().Equal(float64(0.75), output[1][1]) // first average correct
	suite.Assert().Equal(float64(0.0), output[0][1])  // second average correct
}

func (suite *joshTestSuite) TestBGettingLatest() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Query()
		company := u.Get("company")
		suite.Assert().Equal("Test", company)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(2304)
	}))
	result, err := getLatest("Test", ts.URL)
	suite.Assert().Nil(err)
	suite.Assert().Equal(int64(2304), result)
	ts.Close()
}

func (suite *joshTestSuite) TestCGettingRawSentiments() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Query()
		company := u.Get("company")
		before := u.Get("before")
		after := u.Get("after")
		suite.Assert().Equal("Test", company)
		suite.Assert().Equal("7600", before)
		suite.Assert().Equal("3000", after)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([][]float64{[]float64{2000, 0, 0}, []float64{2001, 1, 1}})
	}))

	result, err := getrawSentiments("Test", ts.URL, time.Unix(7600, 0), 6600, 3600)
	suite.Assert().Nil(err)
	suite.Assert().Equal(2, len(result))
	suite.Assert().Equal(3, len(result[0]))
	suite.Assert().Equal(3, len(result[1]))
	suite.Assert().Equal(float64(2000), result[0][0])
	suite.Assert().Equal(float64(2001), result[1][0])
	ts.Close()
}

func (suite *joshTestSuite) TestDPostingAverages() {
	// ensure environmental var is not already set
	old_web_key, ok := os.LookupEnv("WEB_KEY")
	if !ok {
		defer os.Unsetenv("WEB_KEY")
	} else {
		defer os.Setenv("WEB_KEY", old_web_key)
	}
	os.Setenv("WEB_KEY", "Test Web Key")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.Assert().Equal("Test Web Key", r.Header.Get("webkey"))
		type windowAvePayload struct {
			Company  string    `json:"company"`
			Unix     []int     `json:"unix"`
			Averages []float64 `json:"average"`
		}
		payload := &windowAvePayload{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(payload)
		suite.Assert().Nil(err)
		suite.Assert().Equal("Test", payload.Company)
		suite.Assert().Equal(2, len(payload.Unix))
		suite.Assert().Equal(2, len(payload.Averages))
		suite.Assert().Equal(3660, payload.Unix[0])
		suite.Assert().Equal(3600, payload.Unix[1])
		suite.Assert().Equal(-.76, payload.Averages[0])
		suite.Assert().Equal(.76, payload.Averages[1])
		w.WriteHeader(http.StatusOK)
	}))

	err := postAverages("Test", ts.URL, [][]float64{[]float64{3660, -.76}, []float64{3600, .76}})
	suite.Assert().Nil(err)
	ts.Close()
}

func TestJosh(t *testing.T) {
	suite.Run(t, new(joshTestSuite))
}
