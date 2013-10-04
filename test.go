package main

import (
	eventsource "eventsource"
	"log"
	"net/http"
	"time"
	"strconv"
)

func main() {
	es := eventsource.New(nil)
	defer es.Close()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/events", es)
	id := 0
	go func() {
		for {
			es.SendMessage("hello", "asd", strconv.Itoa(id))
			log.Printf("Hello has been sent (consumers: %d)", es.ConsumersCount())
			time.Sleep(2 * time.Second)
			id++
		}
	}()
	log.Print("Open URL http://localhost:8080/ in your browser.")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
