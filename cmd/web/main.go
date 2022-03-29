package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		if err = conn.WriteMessage(msgType, msg); err != nil {
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", socketHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
