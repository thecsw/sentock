package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type elephantTestSuite struct {
	suite.Suite
	myelephant elephant
}

func (suite *elephantTestSuite) SetupSuite() {
	var err error
	db, err = suite.myelephant.connectDB("postgresql://sandy:pass@127.0.0.1:5432/sentock?sslmode=disable")
	suite.Assert().Nil(err)
	suite.Assert().NotNil(db)
	suite.Assert().Nil(suite.myelephant.autoMigrate())
}

func (suite *elephantTestSuite) TestACreateCompanies() {
	a, err := suite.myelephant.createCompany("Company A")
	suite.Assert().NotNil(a)
	suite.Assert().Nil(err)
	b, err := suite.myelephant.createCompany("Company B")
	suite.Assert().NotNil(b)
	suite.Assert().Nil(err)
	c, err := suite.myelephant.createCompany("Company C")
	suite.Assert().NotNil(c)
	suite.Assert().Nil(err)
	d, err := suite.myelephant.createCompany("Company D")
	suite.Assert().NotNil(d)
	suite.Assert().Nil(err)
	e, err := suite.myelephant.createCompany("Company E")
	suite.Assert().NotNil(e)
	suite.Assert().Nil(err)
	f, err := suite.myelephant.createCompany("Company F")
	suite.Assert().NotNil(f)
	suite.Assert().Nil(err)
}

func (suite *elephantTestSuite) TestBCreateSentiments() {
	sa1, err := suite.myelephant.createSentiment("twa1", 1, 1, "Company A")
	suite.Assert().NotNil(sa1)
	suite.Assert().Nil(err)
	sa2, err := suite.myelephant.createSentiment("twa2", 2, 2, "Company A")
	suite.Assert().NotNil(sa2)
	suite.Assert().Nil(err)
	sa3, err := suite.myelephant.createSentiment("twa3", 3, 3, "Company A")
	suite.Assert().NotNil(sa3)
	suite.Assert().Nil(err)

	sb1, err := suite.myelephant.createSentiment("twb1", 1, 1, "Company B")
	suite.Assert().NotNil(sb1)
	suite.Assert().Nil(err)
	sb2, err := suite.myelephant.createSentiment("twb2", 2, 2, "Company B")
	suite.Assert().NotNil(sb2)
	suite.Assert().Nil(err)
	sb3, err := suite.myelephant.createSentiment("twb3", 3, 3, "Company B")
	suite.Assert().NotNil(sb3)
	suite.Assert().Nil(err)

	sc1, err := suite.myelephant.createSentiment("twc1", 1, 1, "Company C")
	suite.Assert().NotNil(sc1)
	suite.Assert().Nil(err)
	sc2, err := suite.myelephant.createSentiment("twc2", 2, 2, "Company C")
	suite.Assert().NotNil(sc2)
	suite.Assert().Nil(err)

	sd1, err := suite.myelephant.createSentiment("twd1", 1, 1, "Company D")
	suite.Assert().NotNil(sd1)
	suite.Assert().Nil(err)
}

func (suite *elephantTestSuite) TestCGetSentiments() {
	sntsa, err := suite.myelephant.getSentiments("Company A", 20, 0)
	suite.Assert().Nil(err)
	suite.Assert().Equal(3, len(sntsa))
	suite.Assert().Equal("twa3", sntsa[0].TweetID)
	suite.Assert().Equal("twa2", sntsa[1].TweetID)
	suite.Assert().Equal("twa1", sntsa[2].TweetID)

	sntsb, err := suite.myelephant.getSentiments("Company B", 20, 1)
	suite.Assert().Nil(err)
	suite.Assert().Equal(2, len(sntsb))

	sntsc, err := suite.myelephant.getSentiments("Company C", 20, 0)
	suite.Assert().Nil(err)
	suite.Assert().Equal(2, len(sntsc))

	sntsd, err := suite.myelephant.getSentiments("Company D", 20, 0)
	suite.Assert().Nil(err)
	suite.Assert().Equal(1, len(sntsd))
}

func TestElephant(t *testing.T) {
	suite.Run(t, new(elephantTestSuite))
}

func (suite *elephantTestSuite) TearDownSuite() {
	suite.Assert().Nil(suite.myelephant.close())
}
