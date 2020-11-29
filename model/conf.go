package model

import (
	"fmt"
	"net/url"
)

type DbConf struct {
	username string
	password string
	scheme   string
	host     string
	port     int
	database string
	flags    map[string][]string
}

func NewDbConf(username string, password string, scheme string, host string, port int, database string) *DbConf {
	return &DbConf{
		username: username,
		password: password,
		scheme:   scheme,
		host:     host,
		port:     port,
		database: database,
		flags:    make(map[string][]string),
	}
}

func (c *DbConf) Flag(name string, value ...string) {
	c.flags[name] = value
}

func (c *DbConf) Dsn() string {
	dsn := url.URL{
		User:     url.UserPassword(c.username, c.password),
		Scheme:   c.scheme,
		Host:     fmt.Sprintf("%s:%d", c.host, c.port),
		Path:     c.database,
		RawQuery: url.Values(c.flags).Encode(),
	}
	return dsn.String()
}
