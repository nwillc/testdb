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
	EmbeddedPostgres(suite.T(), conf)
	db, err := gorm.Open(postgres.Open(conf.String()), &gorm.Config{})
	if err != nil {
		suite.T().Error(err)
	}
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
	if err := suite.db.AutoMigrate(&Person{}); err != nil {
		suite.T().Error(err)
	}

	p1 := Person{FirstName: "John", LastName: "Doe"}

	suite.db.Create(&p1)
	var p2 Person
	suite.db.Last(&p2)
	assert.Equal(suite.T(), p1.ID, p2.ID)
	assert.Equal(suite.T(), p1.FirstName, p2.FirstName)
	assert.Equal(suite.T(), p1.LastName, p2.LastName)
}
