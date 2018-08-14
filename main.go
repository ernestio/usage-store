/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/nats-io/go-nats"
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

	handlers := map[string]nats.MsgHandler{
		"usage.get":  handler.Get,
		"usage.del":  handler.Del,
		"usage.set":  handler.Set,
		"usage.find": handler.Find,
	}
	for subject, h := range handlers {
		if _, err = n.Subscribe(subject, h); err != nil {
			log.Println("Error subscribing " + subject)
		}
	}

	trackables := []string{"instance", "virtual_machine"}
	handlers = map[string]nats.MsgHandler{
		".create.*.done": AddTrackableHandler,
		".delete.*.done": RmTrackableHandler,
		".update.*.done": UpdateTrackableHandler,
	}

	for _, t := range trackables {
		log.Println("Listening for " + t + ".*.*.*")
		for subject, h := range handlers {
			if _, err = n.Subscribe(t+subject, h); err != nil {
				log.Println("Error subscribing "+t+subject, t)
			}
		}
	}
}

func main() {
	setupNats()
	setupPg()
	startHandler()

	runtime.Goexit()
}
