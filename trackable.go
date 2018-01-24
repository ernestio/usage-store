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
}

func addTrackable(msg *nats.Msg) {
	i := trackable{}
	now := time.Now().Unix()
	if err = json.Unmarshal(msg.Data, &i); err != nil {
		log.Println(err)
		return
	}
	t := strings.Split(msg.Subject, ".")
	e := Entity{
		Service: i.Service,
		Name:    i.Name,
		Type:    t[0],
		From:    now,
	}
	if err := e.Save(); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Added trackable " + t[1] + ":" + i.Name)
	}
}

func rmTrackable(msg *nats.Msg) {
	i := trackable{}
	if err = json.Unmarshal(msg.Data, &i); err != nil {
		log.Println(err)
		return
	}
	t := strings.Split(msg.Subject, ".")

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
			if err := entity.Save(); err != nil {
				log.Println(err.Error())
			}
			return
		}
	}
}
