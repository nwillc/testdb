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
 *
 *
 */

package model

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"godb/dbutil"
	"testing"
)

type PersonGoPGTestSuite struct {
	suite.Suite
	db *pg.DB
}

func (suite *PersonGoPGTestSuite) SetupTest() {
	conf := dbutil.NewDbConf(
		"test",
		"test",
		"postgres",
		5432,
		"test",
	)
	// Fire up the embedded Postgres
	EmbeddedPostgres(suite.T(), &dbutil.PostgresContainerConf{DbConf: conf, Image: "postgres:12.4-alpine"})
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.MappedHost(), conf.MappedPort()),
		User:     conf.Username(),
		Password: conf.Password(),
		Database: conf.Database(),
	})
	suite.T().Cleanup(func() {
		_ = db.Close()
	})
	suite.db = db
	err := suite.db.Model(&Person{}).CreateTable(&orm.CreateTableOptions{Temp: true})
	assert.NoError(suite.T(), err)
}

func TestPersonGoPGTestSuite(t *testing.T) {
	suite.Run(t, new(PersonGoPGTestSuite))
}

func (suite *PersonGoPGTestSuite) TestWrite() {
	// Create a person
	p1 := Person{FirstName: "John", LastName: "Doe"}

	// Persist it to database
	_, err := suite.db.Model(&p1).Insert()
	assert.NoError(suite.T(), err)

	// Select all
	var people []Person
	err = suite.db.Model(&people).Select()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), people, 1)
	assert.Equal(suite.T(), p1, people[0])
}
