package main

import (
	eventsource "eventsource"
	"log"
	"net/http"
	"time"
)

func main() {

	es := eventsource.New(nil)
	defer es.Close()
	http.Handle("/", es)

	go func() {
		for {
			es.ProcessMessages()
			time.Sleep( 1 * time.Second )
		}
	}()

	log.Print("Running.")

	http.ListenAndServe(":5151", nil)
}
