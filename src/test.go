package main

import (
	. "google/appengine/channel"
	"log"
)

func main() {
	log.Println("started")
	stop_chan := make(chan bool)

	channel := NewChannel("http://app.myalert.info/online_get_token?id=", "pi")
	socket := channel.Open()
	socket.OnOpened = func() {
		log.Println("socket opened!")
	}

	socket.OnClose = func() {
		log.Println("socket closed!")
		stop_chan <- true
	}

	socket.OnMessage = func(msg *Element) {
		log.Println(msg.ToString())
	}

	socket.OnError = func(err error) {
		log.Println("error:", err)
	}

	<-stop_chan
}
