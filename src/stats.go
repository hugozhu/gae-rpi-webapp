package main

import (
	"analytics"
	"encoding/json"
	. "github.com/hugozhu/gae-channel"
	"log"
	"strings"
	"time"
)

var stats = analytics.NewAnalytics(5, 60)

func init() {

}

func main() {
	for {
		run()
	}
}

func run() {
	log.Println("started")
	stop_chan := make(chan bool)

	channel := NewChannel("http://app.myalert.info/online_get_token?id=pi")
	socket := channel.Open()
	socket.OnOpened = func() {
		log.Println("socket opened!")
	}

	socket.OnClose = func() {
		log.Println("socket closed!")
		stop_chan <- true
	}

	socket.OnMessage = func(msg *Message) {
		if msg.Level() >= 3 && msg.Child.Key == "c" {
			v1 := *msg.Child.Child.Val
			if len(v1) > 0 {
				s := "[" + v1[0].Key + "]"
				var v []string
				json.Unmarshal([]byte(s), &v)
				if len(v) == 2 && v[0] == "ae" {
					s = v[1]
					v = strings.Split(s, "\n")
					zcookie := v[0]
					t := time.Now()
					stats.AddOne(zcookie, t)
				}
			}
		}
	}

	socket.OnError = func(err error) {
		log.Println("error:", err)
	}

	go func() {
		for {
			a, slots_a, b, slots_b := stats.Sum()
			log.Println("UV", a, slots_a)
			log.Println("PV", b, slots_b)
			time.Sleep(5 * time.Second)
		}
	}()

	<-stop_chan
}
