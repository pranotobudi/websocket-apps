package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// Read messages from socket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}
		log.Printf("msg: %s", string(msg))
	}
}

func main() {
	http.HandleFunc("/", ws)
	log.Println("server run on port :8000..")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
