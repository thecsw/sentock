package main

import "github.com/jinzhu/gorm"

var (
	db     *gorm.DB
	dbType = "postgres"
	//connectionString = "postgresql://sandy:pass@db:5432/sentock?sslmode=disable"
	connectionString = "postgresql://sandy:pass@postgres:5432/sentock?sslmode=disable"

	cachedNames map[string]uint = make(map[string]uint)
)

type Company struct {
	gorm.Model `json:"-"`

	Name       string      `gorm:"unique" json:"name"`
	Sentiments []Sentiment `json:"-"`
}

type Sentiment struct {
	gorm.Model `json:"-"`

	TweetID   string  `gorm:"unique;" json:"tweet_id"`
	Unix      int     `json:"unix"`
	Sentiment float64 `json:"sentiment"`
	CompanyID uint    `json:"-"`
}

type Average struct {
	gorm.Model `json:"-"`

	Unix      int     `json:"unix"`
	Average   float64 `json:"average"`
	CompanyID uint    `json:"-"`
}