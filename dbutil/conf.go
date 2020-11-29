package dbutil

import (
	"fmt"
	"net/url"
)

type DbConf struct {
	username   string
	password   string
	scheme     string
	port       int
	mappedHost string
	mappedPort int
	database   string
	flags      map[string][]string
}

//DbConf implements fmt.Stringer
var _ fmt.Stringer = (*DbConf)(nil)

func (c *DbConf) Port() int {
	return c.port
}

func (c *DbConf) Database() string {
	return c.database
}

func (c *DbConf) Username() string {
	return c.username
}

func (c *DbConf) Password() string {
	return c.password
}

func NewDbConf(username string, password string, scheme string, port int, database string) *DbConf {
	return &DbConf{
		username: username,
		password: password,
		scheme:   scheme,
		port:     port,
		database: database,
		flags:    make(map[string][]string),
	}
}

func (c *DbConf) Flag(name string, value ...string) *DbConf {
	c.flags[name] = value
	return c
}

func (c *DbConf) Mapped(mappedHost string, mappedPort int) *DbConf {
	c.mappedHost = mappedHost
	c.mappedPort = mappedPort
	return c
}

func (c *DbConf) String() string {
	dsn := url.URL{
		User:     url.UserPassword(c.username, c.password),
		Scheme:   c.scheme,
		Host:     fmt.Sprintf("%s:%d", c.mappedHost, c.mappedPort),
		Path:     c.database,
		RawQuery: url.Values(c.flags).Encode(),
	}
	return dsn.String()

}
