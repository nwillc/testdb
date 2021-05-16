package dbutil

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq" // Make Postgres lib available for sql.Open
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

// Implement interface
var _ wait.Strategy = (*PostgresStrategy)(nil)

// PostgresStrategy a wait.Strategy to wait on Postgres to start.
type PostgresStrategy struct {
	Port           nat.Port
	startupTimeout time.Duration
	dbConf         *PostgresContainerConf
}

// NewPostgresStrategy constructs a default host port strategy
func NewPostgresStrategy(port nat.Port, dbConf *PostgresContainerConf) *PostgresStrategy {
	return &PostgresStrategy{
		Port:           port,
		startupTimeout: 60 * time.Second,
		dbConf:         dbConf,
	}
}

// WithStartupTimeout sets startupTimeout
func (hp *PostgresStrategy) WithStartupTimeout(startupTimeout time.Duration) *PostgresStrategy {
	hp.startupTimeout = startupTimeout
	return hp
}

// WaitUntilReady implements Strategy.WaitUntilReady
func (hp *PostgresStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) (err error) {
	// limit context to startupTimeout
	ctx, cancelContext := context.WithTimeout(ctx, hp.startupTimeout)
	defer cancelContext()

	var waitInterval = 100 * time.Millisecond

	var port nat.Port
	var i = 0
	for port == "" {
		i++
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s:%w", ctx.Err(), err)
		case <-time.After(waitInterval):
			port, err = mappedPort(ctx, target, hp.Port)
			if err != nil {
				fmt.Printf("(%d) [%s] %s\n", i, port, err)
			}
		}
	}

	psqlInfo := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
		port.Int(), hp.dbConf.Username(), hp.dbConf.Password(), hp.dbConf.Database())

	var success bool
	for !success {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s:%w", ctx.Err(), err)
		case <-time.After(waitInterval):
			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				continue
			}
			_, err = db.ExecContext(ctx, "SELECT 1")
			_ = db.Close()
			if err == nil {
				success = true
			}
		}
	}
	return nil
}

func mappedPort(ctx context.Context, target wait.StrategyTarget, port nat.Port) (nat.Port, error) {
	var rp nat.Port
	var rerr = fmt.Errorf("failed to get mapped port")
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	rp, rerr = target.MappedPort(ctx, port)
	return rp, rerr
}
