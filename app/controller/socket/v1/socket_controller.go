package controllerSocketV1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Define our message object
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}


func HandleConnection(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	clients[conn] = true

	go func(conn *websocket.Conn) {

		for {

			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				conn.Close()
				delete(clients, conn)
				break
			}

			// Send it out to every client that is currently connected
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}(conn)
}

