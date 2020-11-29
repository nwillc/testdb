package model

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"godb/dbutil"
	"testing"
	"time"
)

const (
	image    = "postgres:latest"
	logMsg   = "database system is ready to accept connections"
)

func EmbeddedPostgres(t *testing.T, conf *dbutil.DbConf) {
	ctx := context.Background()
	natPort := fmt.Sprintf("%d/tcp", conf.Port())
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{ natPort },
		Env: map[string]string{
			"POSTGRES_PASSWORD": conf.Password(),
			"POSTGRES_USER":     conf.Username(),
			"POSTGRES_DATABASE": conf.Database(),
		},
		WaitingFor: wait.ForLog(logMsg).
			WithPollInterval(100 * time.Millisecond).
			WithOccurrence(2),
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	// Even after log message found Postgres needs a touch more...
	time.Sleep(200 * time.Millisecond)
	mp, err := pg.MappedPort(ctx, nat.Port(natPort))
	if err != nil {
		t.Error(err)
	}
	ma, err := pg.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	// When test is done terminate container
	t.Cleanup(func() {
		_ = pg.Terminate(ctx)
	})
	// Note the containers mapped host and port
	conf.Mapped(ma, mp.Int())
}
