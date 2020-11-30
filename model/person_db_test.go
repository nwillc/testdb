package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"godb/dbutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type PersonDbTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *PersonDbTestSuite) SetupTest() {
	conf := dbutil.NewDbConf(
		"postgres",
		"admin",
		"postgres",
		5432,
		"postgres",
	)
	// Fire up the embedded Postgres
	EmbeddedPostgres(suite.T(), conf)
	// Open the gorm connection to it
	db, err := gorm.Open(
		postgres.Open(conf.String()),
		&gorm.Config{},
	)
	if err != nil {
		suite.T().Error(err)
	}
	// Explicitly clean up gorm after the test
	suite.T().Cleanup(func() {
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()
	})
	suite.db = db
}

func TestPersonDbTestSuite(t *testing.T) {
	suite.Run(t, new(PersonDbTestSuite))
}

func (suite *PersonDbTestSuite) TestWrite() {
	// Use gorm's migration to set up our table
	if err := suite.db.AutoMigrate(&Person{}); err != nil {
		suite.T().Error(err)
	}

	// Create a person
	p1 := Person{FirstName: "John", LastName: "Doe"}

	// Persist it to database
	suite.db.Create(&p1)
	var p2 Person

	// Find the last Person in the database
	suite.db.Last(&p2)

	// Compare...
	assert.Equal(suite.T(), p1.ID, p2.ID)
	assert.Equal(suite.T(), p1.FirstName, p2.FirstName)
	assert.Equal(suite.T(), p1.LastName, p2.LastName)
}
