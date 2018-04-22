package main

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"fmt"
)

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func Echo(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Print("Cat't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)
		msg := reply
		fmt.Println("Sending to client: " + msg)

		log.Println(ws.Config())
		//log.Println(ws.LocalAddr())

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Cat't send")
			break
		}
	}
}
