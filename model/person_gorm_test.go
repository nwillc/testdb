/*
 * Copyright (c) 2020, nwillc@gmail.com
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

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
		"test",
		"test",
		"postgres",
		5432,
		"test",
	)
	// Fire up the embedded Postgres
	EmbeddedPostgres(suite.T(), &dbutil.PostgresContainerConf{DbConf: conf, Image: "postgres:12.4-alpine"})
	// Open the gorm connection to it
	db, err := gorm.Open(postgres.Open(conf.String()), &gorm.Config{})
	if err != nil {
		suite.T().Error(err)
	}
	// Explicitly clean up gorm after the test
	suite.T().Cleanup(func() {
		sqlDb, _ := db.DB()
		_ = sqlDb.Close()
	})
	suite.db = db
	// Use gorm's migration to set up our table
	if err := suite.db.AutoMigrate(&Person{}); err != nil {
		assert.NoError(suite.T(), err)
	}
}

func TestPersonDbTestSuite(t *testing.T) {
	suite.Run(t, new(PersonDbTestSuite))
}

func (suite *PersonDbTestSuite) TestWrite() {
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
