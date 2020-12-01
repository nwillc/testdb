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
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"godb/dbutil"
	"testing"
)

type PostgresContainerConf struct {
	*dbutil.DbConf
	Image string
}

//EmbeddedPostgres spins up a Postgres container.
func EmbeddedPostgres(t *testing.T, conf *PostgresContainerConf) {
	t.Helper()
	ctx := context.Background()
	natPort := fmt.Sprintf("%d/tcp", conf.Port())
	// Configure the container
	req := testcontainers.ContainerRequest{
		Image:        conf.Image,
		ExposedPorts: []string{natPort},
		Env: map[string]string{
			"POSTGRES_PASSWORD": conf.Password(),
			"POSTGRES_USER":     conf.Username(),
			"POSTGRES_DB":       conf.Database(),
		},
		WaitingFor: wait.ForListeningPort(nat.Port(natPort)),
	}
	// Spin up the container
	pg, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		t.Error(err)
	}
	// Even after log message found Postgres needs a touch more...
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
