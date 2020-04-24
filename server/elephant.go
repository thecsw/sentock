package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	// We need this for gorm to connect to postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type elephant struct{}

var (
	Elephant = &elephant{}
)

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

func (*elephant) createCompany(name string) (*Company, error) {
	res := &Company{Name: name}
	return res, db.Create(res).Error
}

func (e *elephant) createSentiment(tweetID string, unix int, sentiment float64, company string) (*Sentiment, error) {
	ans := &Sentiment{
		TweetID:   tweetID,
		Unix:      unix,
		Sentiment: sentiment,
		CompanyID: e.ctod(company, true),
	}
	return ans, db.Create(ans).Error
}

func (e *elephant) createWindowAverage(unix []int, averages []float64, company string) ([]*Average, error) {
	end := make([]*Average, len(unix))
	var err error
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

func (e *elephant) getSentiments(company string, before, after int) ([]Sentiment, error) {
	result := make([]Sentiment, 0, 128)
	return result, db.Model(e.fcompany(company)).
		Where("unix > ?", after).
		Where("unix < ?", before).
		Order("unix desc").
		Related(&result).
		Error
}

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
	}
	return result.Unix, err
}

func (e *elephant) getAverages(company string, before, after int) ([]Average, error) {
	result := make([]Average, 0, 128)
	return result, db.Model(e.fcompany(company)).
		Where("unix > ?", after).
		Where("unix < ?", before).
		Order("unix desc").
		Related(&result).
		Error
}

func (*elephant) close() error {
	return db.Close()
}

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
