package main

import (
	"net/http"
	"code.google.com/p/go.net/websocket"
)

func reportEvent( ws *websocket.Conn ) {
}

func main() {
	http.Handle("/reportEvent", websocket.Handler( reportEvent ) )
	http.ListenAndServe( ":5152", nil )
}
