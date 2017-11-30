package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats"
)

type trackable struct {
	Service string `json:"service"`
	Name    string `json:"name"`
	Powered *bool  `json:"powered"`
}

// UpdateTrackableHandler : handler to manage trackable modifications
func UpdateTrackableHandler(msg *nats.Msg) {
	i := trackable{}
	if err = json.Unmarshal(msg.Data, &i); err != nil {
		log.Println(err)
		return
	}
	if i.Powered == nil {
		return
	} else if *i.Powered == true {
		addTrackable(msg.Subject, i)
	} else {
		rmTrackable(msg.Subject, i)
	}
}

// AddTrackableHandler : Handler to manage trackable additions
func AddTrackableHandler(msg *nats.Msg) {
	i := trackable{}
	if err = json.Unmarshal(msg.Data, &i); err != nil {
		log.Println(err)
		return
	}
	addTrackable(msg.Subject, i)
}

// RmTrackableHandler : Handler to manage trackable removals
func RmTrackableHandler(msg *nats.Msg) {
	i := trackable{}
	if err = json.Unmarshal(msg.Data, &i); err != nil {
		log.Println(err)
		return
	}
	rmTrackable(msg.Subject, i)
}

func addTrackable(subject string, i trackable) {
	now := time.Now().Unix()
	t := strings.Split(subject, ".")
	e := Entity{
		Service: i.Service,
		Name:    i.Name,
		Type:    t[0],
		From:    now,
	}
	e.Save()
	log.Println("Added trackable " + t[1] + ":" + i.Name)
}

func rmTrackable(subject string, i trackable) {
	t := strings.Split(subject, ".")

	e := Entity{
		Service: i.Service,
		Name:    i.Name,
		Type:    t[0],
	}

	for _, v := range e.Find() {
		entity := v.(Entity)
		if entity.To == 0 {
			now := time.Now().Unix()
			entity.To = now
			entity.Save()
			return
		}
	}
}
