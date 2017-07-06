package main

import (
	"log"
	"net/http"

	"time"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type Message struct {
	Time    string `json:"time"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

type chatroom struct {
	in      chan Message     // incoming messages
	join    chan *client     // handles new user notigications
	leave   chan *client     // handles users leaving
	clients map[*client]bool //canonical sourcce for who is in the room
	history []Message        // old messages
}

// makes creation of a chatroom object less painful
func newChatRoom() *chatroom {
	return &chatroom{
		in:      make(chan Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// this method maintains the chatroom's state
func (cr *chatroom) run() {
	for {
		select {
		case client := <-cr.join:
			// keep track of this new user and send it the preexisting messages
			cr.clients[client] = true
			for _, msg := range cr.history {
				client.out <- msg
			}
		case client := <-cr.leave:
			// remove user from list and close its outgoing message channel
			delete(cr.clients, client)
			close(client.out)
		case msg := <-cr.in:
			// send new message to usrs and store it for later use
			for client := range cr.clients {
				client.out <- msg
			}
			// @todo: potentially sore this in redis or something similar
			cr.history = append(cr.history, msg)
		}
	}
}

// used for altering the new websocket connections
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// ServeHTTP turns this type into an instance of the Handler interface and
// orchestrates updates to the user
func (cr *chatroom) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := &client{
		name:      req.RemoteAddr,
		websocket: socket,
		out:       make(chan Message),
		chatroom:  cr,
	}

	cr.join <- client
	defer func() { cr.leave <- client }()
	go client.Write()
	client.Read()
}

// keeps track of individual users
type client struct {
	name      string
	websocket *websocket.Conn
	out       chan Message
	chatroom  *chatroom
}

// Read repeatedly checks the user's socket for new messages from the user's browser
func (c *client) Read() {
	defer c.websocket.Close()
	for {
		var msg Message
		err := c.websocket.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message", err)
			break
		}

		t := time.Now()

		msg.ID = c.name
		msg.Time = t.Format(time.Kitchen)

		c.chatroom.in <- msg
	}
}

// Write sends outgoing messages from other users to this user's browser
func (c *client) Write() {
	defer c.websocket.Close()
	for m := range c.out {
		err := c.websocket.WriteJSON(m)
		if err != nil {
			log.Println("Error writing message", err)
			break
		}
	}
}
