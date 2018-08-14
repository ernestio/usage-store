/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/r3labs/natsdb"
)

// Entity : the database mapped entity
type Entity struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Service   string `json:"service"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	From      int64  `json:"from" gorm:"column:from_date"`
	To        int64  `json:"to" gorm:"column:to_date"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// TableName : set Entity's table name to be usage
func (Entity) TableName() string {
	return "usage"
}

// Find : based on the defined fields for the current entity
// will perform a search on the database
func (e *Entity) Find() (list []interface{}) {
	entities := []Entity{}
	if e.From != 0 && e.To != 0 {
		log.Println("Querying from & to")
		db.Where("from_date <= ? OR to_date >= ?", e.To, e.From).Find(&entities)
	} else if e.From != 0 {
		log.Println("Querying from")
		db.Where("to_date >= ?", e.From).Find(&entities)
	} else if e.To != 0 {
		println(e.To)
		log.Println("Querying to")
		db.Where("from_date < ?", e.To).Find(&entities)
	} else {
		db.Find(&entities)
	}
	list = make([]interface{}, len(entities))
	for i, s := range entities {
		list[i] = s
	}

	return list
}

// MapInput : maps the input []byte on the current entity
func (e *Entity) MapInput(body []byte) {
	if err := json.Unmarshal(body, &e); err != nil {
		log.Println("Invalid input " + err.Error())
	}
}

// HasID : determines if the current entity has an id or not
func (e *Entity) HasID() bool {
	if e.ID == 0 {
		return false
	}
	return true
}

// LoadFromInput : Will load from a []byte input the database stored entity
func (e *Entity) LoadFromInput(msg []byte) bool {
	e.MapInput(msg)
	var stored Entity
	if e.ID != 0 {
		db.First(&stored, e.ID)
	} else if e.Name != "" {
		db.Where("name = ?", e.Name).First(&stored)
	}
	if &stored == nil {
		return false
	}
	if ok := stored.HasID(); !ok {
		return false
	}

	e.ID = stored.ID
	e.Service = stored.Service
	e.Name = stored.Name
	e.Type = stored.Type
	e.From = stored.From
	e.To = stored.To
	e.CreatedAt = stored.CreatedAt
	e.UpdatedAt = stored.UpdatedAt

	return true
}

// LoadFromInputOrFail : Will try to load from the input an existing entity,
// or will call the handler to Fail the nats message
func (e *Entity) LoadFromInputOrFail(msg *nats.Msg, h *natsdb.Handler) bool {
	stored := &Entity{}
	ok := stored.LoadFromInput(msg.Data)
	if !ok {
		h.Fail(msg)
	}
	*e = *stored

	return ok
}

// Update : It will update the current entity with the input []byte
func (e *Entity) Update(body []byte) error {
	e.MapInput(body)
	stored := Entity{}
	db.First(&stored, e.ID)

	stored.To = e.To

	db.Save(&stored)
	e = &stored

	return nil
}

// Delete : Will delete from database the current Entity
func (e *Entity) Delete() error {
	db.Unscoped().Delete(&e)

	return nil
}

// Save : Persists current entity on database
func (e *Entity) Save() error {
	db.Save(&e)

	return nil
}
