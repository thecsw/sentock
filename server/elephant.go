package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	// We need this for gorm to connect to postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func connectDB(database string) (*gorm.DB, error) {
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

func createCompany(name string) error {
	return db.Create(&Company{Name: name}).Error
}

func createSentiment(tweetID string, unix int, sentiment float64, company string) error {
	return db.Create(&Sentiment{
		TweetID:   tweetID,
		Unix:      unix,
		Sentiment: sentiment,
		CompanyID: ctod(company),
	}).Error
}

func getSentiments(company string, before, after int) ([]Sentiment, error) {
	result := make([]Sentiment, 0, 128)
	return result, db.Model(&Sentiment{}).
		Where("id = ?", ctod(company)).
		Where("unix > ?", after).
		Where("unix < ?", before).
		Find(&result).
		Error
}

func close() error {
	return db.Close()
}

func autoMigrate() error {
	return db.AutoMigrate(&Company{}, &Sentiment{}).Error
}

func ctod(company string) uint {
	companyID := uint(0)
	var ok bool
	if companyID, ok = cachedNames[company]; !ok {
		res := &Company{}
		err := db.Where(&Company{Name: company}).First(res).Error
		if err != nil {
			return 0
		}
		cachedNames[company] = res.ID
		companyID = res.ID
	}
	return companyID
}