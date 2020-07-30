package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

func main() {
	app := fiber.New()
	count := 0

	var mtx sync.Mutex

	app.Use(func(c *fiber.Ctx) {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			c.Next()
		}
	})

	app.Get("/ws/groupchat", websocket.New(func(c *websocket.Conn) {
		c.Locals("allowed")

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				// log.Println("read:", err)
				break
			}
			// log.Printf("recv: %s", msg)
			mtx.Lock()
			count++
			mtx.Unlock()
			if count%20 == 0 {
				log.Println(count)
			}

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	log.Fatal(app.Listen(8000))
}
