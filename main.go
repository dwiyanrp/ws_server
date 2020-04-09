package main

import (
	"net/http"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var hub = NewHub()
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 10 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	hub.ActivateChannel("2")
	go hub.HandleMessages()
}

func main() {
	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	r := gin.New()

	r.GET("/ws/groupchat", HandleGroupchat)
	r.GET("/channel/:channel_id", HandleGetChannel)
	r.GET("/channel/:channel_id/activate", HandleActivateChannel)
	r.GET("/channel/:channel_id/deactivate", HandleDeactivateChannel)

	r.Run(":8000")
}
