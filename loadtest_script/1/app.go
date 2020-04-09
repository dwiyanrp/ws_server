package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	start := 1
	for i := start; i <= 4000-start; i++ {
		go createClient(i)
		log.Println(i)

		if i%10 == 0 {
			time.Sleep(50 * time.Millisecond)
		} else {
			time.Sleep(3 * time.Millisecond)
		}
	}

	time.Sleep(time.Hour)
}

func createClient(id int) {
	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws/groupchat", RawQuery: fmt.Sprintf("channel_id=2&user_id=%d&device=%d", id, id)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
	} else {
		tick := time.NewTicker(20 * time.Second)
		for {
			select {
			case <-tick.C:
				c.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}
}
