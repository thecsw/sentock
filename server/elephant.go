package main

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	// We need this for gorm to connect to postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// elephant is a struct to define DB interaction
type elephant struct{}

var (
	// Elephant is the exposed model into DB
	Elephant = &elephant{}
)

// connectDB connects to the database
func (*elephant) connectDB(database string) (*gorm.DB, error) {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	tries := 0
	for tries = 0; time.Now().Before(deadline); tries++ {
		db, err := gorm.Open(dbType, database)
		if err == nil {
			return db, nil
		}
		log.Printf("Couldn't connect to the database (%s). Retrying...", err.Error())
		time.Sleep(time.Second << uint(tries))
	}
	// Sleep for 2 seconds
	time.Sleep(2 * time.Second)
	return nil, fmt.Errorf("failed connect to the database after %d attempts", tries)
}

// createCompany adds company to the database
func (*elephant) createCompany(name string) (*Company, error) {
	res := &Company{Name: name}
	return res, db.Create(res).Error
}

// createSentiment adds a raw sentiment datapoint to the database
func (e *elephant) createSentiment(tweetID string, unix int, sentiment float64, company string) (*Sentiment, error) {
	ans := &Sentiment{
		TweetID:   tweetID,
		Unix:      unix,
		Sentiment: sentiment,
		CompanyID: e.ctod(company, true),
	}
	return ans, db.Create(ans).Error
}

// createWindowAverage adds a list of averaged sentiment datapoints to the database
func (e *elephant) createWindowAverage(unix []int, averages []float64, company string) ([]*Average, error) {
	end := make([]*Average, len(unix))
	var err error
	if len(averages) != len(unix) {
		return nil, errors.New("Length of unix timestamps and averages don't match!")
	}
	for ind, time := range unix {
		ave := &Average{
			Unix:      time,
			Average:   averages[ind],
			CompanyID: e.ctod(company, true),
		}
		end = append(end, ave)
		err = db.Create(ave).Error
		if err != nil {
			break
		}
	}
	return end, err
}

// getSentiments gets raw sentiments from "sentiments" table for all values
// that are connected to company and bound by after < t < before
func (e *elephant) getSentiments(company string, before, after int) ([]Sentiment, error) {
	result := make([]Sentiment, 0, 128)
	return result, db.Model(e.fcompany(company)).
		Where("unix > ?", after).
		Where("unix < ?", before).
		Order("unix desc").
		Related(&result).
		Error
}

// getCompanies returns a string list of all companies in "companies" table
func (e *elephant) getCompanies() ([]string, error) {
	result := make([]string, 0, 16)
	return result, db.Model(&Company{}).
		Pluck("name", &result).
		Error
}

// getLatestAverageSentiment gets latest unix timestamp for
// company. Used by josh for anchoring and starting processing
// only from that point onward
func (e *elephant) getLatestAverageSentiment(company string) (int, error) {
	result := &Average{}
	err := db.Model(e.fcompany(company)).
		Order("unix desc").
		Limit(1).
		Related(result).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return 0, nil
		}
		return 0, err
	}
	return result.Unix, err
}

// getAverages is the same as getSentiments but only returns already
// processed values
func (e *elephant) getAverages(company string, before, after int) ([]Average, error) {
	result := make([]Average, 0, 128)
	return result, db.Model(e.fcompany(company)).
		Where("unix > ?", after).
		Where("unix < ?", before).
		Order("unix asc").
		Related(&result).
		Error
}

// deleteOldRaws deletes all raw sentiments that are older than before
func (e *elephant) deleteOldRaws(before int64) error {
	return db.Unscoped().Model(&Sentiment{}).
		Where("unix < ?", before).
		Delete(&Sentiment{}).
		Error
}

// close closes the connection to the database
func (*elephant) close() error {
	return db.Close()
}

// autoMigrate automatically sets up all the tables that are used
func (*elephant) autoMigrate() error {
	return db.AutoMigrate(&Company{}, &Sentiment{}, &Average{}).Error
}

func (e *elephant) fcompany(company string) *Company {
	res := &Company{}
	res.ID = e.ctod(company, false)
	if res.ID == 0 {
		fmt.Println("failed to resolve company: ", company)
	}
	return res
}

// ctod takes a company and returns the company's id
// (and possibly makes the company a new entry if it wasn't there already)
func (e *elephant) ctod(company string, autocreate bool) uint {
	companyID := uint(0)
	var ok bool
	if companyID, ok = cachedNames[company]; !ok {
		res := &Company{}
		err := db.Where(&Company{Name: company}).First(res).Error
		// Not found, automatically make one
		if res.ID == 0 || err != nil {
			if !autocreate {
				return 0
			}
			c, err := e.createCompany(company)
			if err != nil {
				return 0
			}
			res.ID = c.ID
		}
		cachedNames[company] = res.ID
		companyID = res.ID
	}
	return companyID
}
