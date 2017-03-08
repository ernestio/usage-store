/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats"
	"github.com/r3labs/natsdb"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var n *nats.Conn
var db *gorm.DB
var err error
var handler natsdb.Handler

func startHandler() {
	handler = natsdb.Handler{
		NotFoundErrorMessage:   natsdb.NotFound.Encoded(),
		UnexpectedErrorMessage: natsdb.Unexpected.Encoded(),
		DeletedMessage:         []byte(`{"status":"deleted"}`),
		Nats:                   n,
		NewModel: func() natsdb.Model {
			return &Entity{}
		},
	}

	if _, err = n.Subscribe("usage.get", handler.Get); err != nil {
		log.Println("Error subscribing usage.get")
	}
	if _, err = n.Subscribe("usage.del", handler.Del); err != nil {
		log.Println("Error subscribing usage.del")
	}
	if _, err = n.Subscribe("usage.set", handler.Set); err != nil {
		log.Println("Error subscribing usage.set")
	}
	if _, err = n.Subscribe("usage.find", handler.Find); err != nil {
		log.Println("Error subscribing usage.find")
	}

	// TODO : This should probably be moved to the config service, so we can easily
	// configure it externaly
	trackables := []string{"instance"}
	for _, t := range trackables {
		log.Println("Listening for " + t + ".*.*.*")

		if _, err = n.Subscribe(t+".create.*.done", addTrackable); err != nil {
			log.Println("Error subscribing "+t+".create.*.done", t)
		}
		if _, err = n.Subscribe(t+".delete.*.done", rmTrackable); err != nil {
			log.Println("Error subscribing "+t+".delete.*.done", t)
		}
	}

}

func main() {
	setupNats()
	setupPg()
	startHandler()

	runtime.Goexit()
}
