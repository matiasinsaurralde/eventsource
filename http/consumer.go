package http

import (
	"net"
	"net/http"
	"time"
)

type consumer struct {
	conn   net.Conn
	es     *eventSource
	in     chan []byte
	request	string
	staled bool
}

func newConsumer(resp http.ResponseWriter, es *eventSource, req *http.Request) (*consumer, error) {
	conn, _, err := resp.(http.Hijacker).Hijack()

	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		conn:   conn,
		es:     es,
		in:     make(chan []byte, 10),
		request:	req.URL.RawQuery,
		staled: false,
	}

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\nContent-Type: text/event-stream\nX-Accel-Buffering: no\nAccess-Control-Allow-Methods: GET\nAccess-Control-Allow-Credentials: true\nAccess-Control-Allow-Origin: " + req.Header.Get("Origin") + "\n\n"))
	if err != nil {
		conn.Close()
		return nil, err
	}

	go func() {
		for message := range consumer.in {
			conn.SetWriteDeadline(time.Now().Add(consumer.es.timeout))
			_, err := conn.Write(message)
			if err != nil {
				netErr, ok := err.(net.Error)
				if !ok || !netErr.Timeout() || consumer.es.closeOnTimeout {
					consumer.staled = true
					consumer.conn.Close()
					consumer.es.staled <- consumer
					return
				}
			}
		}
		consumer.conn.Close()
	}()

	return consumer, nil
}
