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

package dbutil

import (
	"fmt"
	"net/url"
)

// DbConf holds the database configuration information.
type DbConf struct {
	username   string
	password   string
	scheme     string
	port       int
	database   string
	mappedHost string
	mappedPort int
	flags      map[string][]string
}

// MappedPort make the mapped port available for read
func (c *DbConf) MappedPort() int {
	return c.mappedPort
}

// MappedHost make the mapped host available fo read
func (c *DbConf) MappedHost() string {
	return c.mappedHost
}

// PostgresContainerConf adds and image name to a DbConf.
type PostgresContainerConf struct {
	*DbConf
	Image string
}

// Database getter
func (c *DbConf) Database() string {
	return c.database
}

// Port getter
func (c *DbConf) Port() int {
	return c.port
}

// Password getter
func (c *DbConf) Password() string {
	return c.password
}

// Username getter
func (c *DbConf) Username() string {
	return c.username
}

// NewDbConf creates a new DbConf with the prerequisite information.
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

// Flag adds a configuration flag to the DbConf
func (c *DbConf) Flag(name string, value ...string) *DbConf {
	c.flags[name] = value
	return c
}

// Mapped adds the host and port mapped by a container.
func (c *DbConf) Mapped(mappedHost string, mappedPort int) *DbConf {
	c.mappedHost = mappedHost
	c.mappedPort = mappedPort
	return c
}

// String implements fmt.Stringer and produces a DSN format string.
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
