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

package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"testdb/dbutil"
	"testdb/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	log.Println("Start")
	const port = 10865
	conf := dbutil.NewDbConf(
		"test",
		"test",
		"postgres",
		port,
		"test",
	).Mapped("localhost", port)
	db, err := gorm.Open(postgres.Open(conf.String()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&model.Person{}); err != nil {
		panic(err)
	}

	p1 := model.Person{FirstName: "John", LastName: "Doe"}
	p2 := model.Person{FirstName: "Jane", LastName: "Smith"}

	db.Create(&p1)
	db.Commit()
	var p3 model.Person
	db.Find(&p3)

	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(p3)
	log.Println("End")
}
