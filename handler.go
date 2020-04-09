package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleGroupchat(c *gin.Context) {
	channelID, userID, device := c.Query("channel_id"), c.Query("user_id"), c.Query("device")
	if channelID == "" {
		c.String(http.StatusBadRequest, "Must have channel_id")
		return
	} else if userID == "" {
		c.String(http.StatusBadRequest, "Must have user_id")
		return
	} else if device == "" {
		c.String(http.StatusBadRequest, "Must have device")
		return
	}

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	if !hub.IsChannelActive(channelID) {
		ws.WriteMessage(websocket.TextMessage, []byte("Channel : "+channelID+" Deactive"))
		return
	}

	user := hub.GetUser(channelID, userID)
	if user != nil {
		user.Conn.WriteMessage(websocket.TextMessage, []byte("Kick User"))
		user.Conn.Close()
		hub.RemoveUser(channelID, userID)
	}

	hub.AddUser(channelID, userID, device, ws)

	for {
		// Use this if message format is JSON
		// var msg Message
		// err := ws.ReadJSON(&msg)
		// if err != nil {
		// 	hub.RemoveUser(channelID, userID)
		// 	break
		// }

		// Use this if message format is string
		_, p, err := ws.ReadMessage()
		if err != nil {
			hub.RemoveUser(channelID, userID)
			break
		}

		// Send the newly received message to the broadcast channel
		hub.broadcast <- &Message{
			ChannelID: channelID,
			Message:   string(p),
		}
	}
}

func HandleGetChannel(c *gin.Context) {
	channelID := c.Param("channel_id")
	if channelID == "" {
		c.String(http.StatusBadRequest, "Must have channel_id")
		return
	}

	if !hub.IsChannelActive(channelID) {
		c.String(http.StatusBadRequest, "Channel "+channelID+" Deactive")
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("%d Active User\n", len(hub.channels[channelID].Users)))
	for _, user := range hub.channels[channelID].Users {
		c.String(http.StatusOK, fmt.Sprintf("User ID : %s, Device : %s\n", user.UserID, user.Device))
	}
}

func HandleActivateChannel(c *gin.Context) {
	channelID := c.Param("channel_id")
	if channelID == "" {
		c.String(http.StatusBadRequest, "Must have channel_id")
		return
	}

	hub.ActivateChannel(channelID)
	c.String(http.StatusOK, "Channel : %s activated", channelID)
}

func HandleDeactivateChannel(c *gin.Context) {
	channelID := c.Param("channel_id")
	if channelID == "" {
		c.String(http.StatusBadRequest, "Must have channel_id")
		return
	}

	hub.DeactivateChannel(channelID)
	c.String(http.StatusOK, "Channel : %s deactivated", channelID)
}
