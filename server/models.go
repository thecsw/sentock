package main

import "github.com/jinzhu/gorm"

var (
	db               *gorm.DB
	dbType           = "postgres"
	connectionString = "postgresql://sandy:pass@postgres:5432/sentock?sslmode=disable"
)

// Company is the table schema for storing the companies
// we are doing the sentimental analysis on
type Company struct {
	gorm.Model `json:"-"`

	Name       string      `gorm:"unique" json:"name"`
	Sentiments []Sentiment `json:"-"`
}

// Sentiment is the table for raw sentiments
type Sentiment struct {
	gorm.Model `json:"-"`

	TweetID   string  `gorm:"unique;" json:"tweet_id"`
	Unix      int     `json:"unix"`
	Sentiment float64 `json:"sentiment"`
	CompanyID uint    `json:"-"`
}

// Average is the table for processed sentiments
type Average struct {
	gorm.Model `json:"-"`

	Unix      int     `json:"unix"`
	Average   float64 `json:"average"`
	CompanyID uint    `json:"-"`
}
