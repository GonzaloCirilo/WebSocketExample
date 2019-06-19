package sockets

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan Message
}
// Read messages per connection handled by a goroutine
func (c *Client) readIncomingMessage(){
	defer func(){
		c.hub.unregister <- c //unregister client from hub
		c.conn.Close()
	}()

	for{
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message,[]byte{'\n'}, []byte{' '}, -1))
		// process incoming message
		m := Message{string(message) }
		c.hub.messages <- m
	}
}

func (c *Client) writeMessage(){
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case mess, ok := <-c.send:
			if !ok{
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
		c.conn.WriteMessage(websocket.TextMessage, []byte(mess.Text))
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request){
	conn, err := websocket.Upgrader{HandshakeTimeout:1024,ReadBufferSize:1024}.Upgrade(
		w,r,nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan Message, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMessage()
	go client.readIncomingMessage()
}