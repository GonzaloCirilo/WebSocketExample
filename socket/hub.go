package socket

import "encoding/json"

type Message struct{
	Text string `json:"text"`
}

type Hub struct {
	clients map[*Client]bool
	messages chan Message
	register chan *Client //register clients
	unregister chan *Client
}

func NewHub() *Hub{
	return &Hub{
		messages:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run(){
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true // map new client
			jsonMessage, _ := json.Marshal(&Message{ "/A new socket has connected."})
			h.send(Message{string(jsonMessage)}, client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients,client)
				close(client.send)
				jsonMessage, _ := json.Marshal(&Message{ "/A socket has disconnected."})
				h.send(Message{string(jsonMessage)}, client)
			}
		case message := <-h.messages:
			for client := range h.clients{
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}

func (h *Hub) send(message Message, ignore *Client){
	for conn := range h.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}
