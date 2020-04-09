package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Message struct
type Message struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
}

// User struct
type User struct {
	sync.Mutex
	UserID string
	Device string
	Conn   *websocket.Conn
}

// Channel struct
type Channel struct {
	ChannelID string
	Admins    []string
	Users     map[string]*User
}

// NewChannel func
func NewChannel(channelID string) *Channel {
	return &Channel{
		ChannelID: channelID,
		Admins:    []string{},
		Users:     make(map[string]*User),
	}
}

// Hub struct
type Hub struct {
	sync.Mutex
	broadcast chan *Message
	channels  map[string]*Channel
}

// NewHub func
func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan *Message),
		channels:  make(map[string]*Channel),
	}
}

// IsChannelActive func
func (h *Hub) IsChannelActive(channelID string) bool {
	_, exists := h.channels[channelID]
	return exists
}

// ActivateChannel func
func (h *Hub) ActivateChannel(channelID string) {
	if h.channels[channelID] == nil {
		h.channels[channelID] = NewChannel(channelID)
	}
}

// DeactivateChannel func
func (h *Hub) DeactivateChannel(channelID string) {
	if h.channels[channelID] != nil {
		for _, user := range h.channels[channelID].Users {
			user.Conn.WriteMessage(websocket.CloseMessage, nil)
			user.Conn.Close()
		}

		delete(h.channels, channelID)
	}
}

// GetUser func
func (h *Hub) GetUser(channelID string, userID string) *User {
	return h.channels[channelID].Users[userID]
}

// AddUser func
func (h *Hub) AddUser(channelID string, userID string, device string, client *websocket.Conn) {
	h.Lock()
	h.channels[channelID].Users[userID] = &User{UserID: userID, Device: device, Conn: client}
	h.Unlock()
}

// RemoveUser func
func (h *Hub) RemoveUser(channelID string, userID string) {
	h.Lock()
	delete(h.channels[channelID].Users, userID)
	h.Unlock()
}

// HandleMessages func
func (h *Hub) HandleMessages() {
	for {
		msg := <-h.broadcast
		hub.BroadcastMessage(msg.ChannelID, msg.Message)
	}
}

// BroadcastMessage func
func (h *Hub) BroadcastMessage(channelID string, msg string) {
	if h.channels[channelID] != nil {
		for userID, user := range h.channels[channelID].Users {
			if err := user.Conn.WriteJSON(msg); err != nil {
				log.Printf("error: %v", err)
				user.Conn.Close()
				h.RemoveUser(channelID, userID)
			}
		}
	}
}
