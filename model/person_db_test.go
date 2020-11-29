package model

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

const (
	user     = "postgres"
	password = "admin"
	database = "postgres"
	port = "5432/tcp"
)

type PersonDbTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *PersonDbTestSuite) SetupTest() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"POSTGRES_PASSWORD": password,
			"POSTGRES_USER":     user,
			"POSTGRES_ATABASE":  database,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithPollInterval(100 * time.Millisecond),
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		suite.T().Error(err)
	}
	time.Sleep(500 * time.Millisecond)
	mappedPort, err := pg.MappedPort(ctx, port)
	if err != nil {
		suite.T().Error(err)
	}
	host, err := pg.Host(ctx)
	if err != nil {
		suite.T().Error(err)
	}
	suite.T().Cleanup(func() {
		_ = pg.Terminate(ctx)
	})
	conf := NewDbConf(
		user,
		password,
		"postgres",
		host,
		mappedPort.Int(),
		database,
	)
	db, err := gorm.Open(postgres.Open(conf.Dsn()), &gorm.Config{})
	if err != nil {
		suite.T().Error(err)
	}
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
