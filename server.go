package sockets

import (
	"flag"
	"net/http"
)

func runServer(port string) error{
	hub := &Hub{
		messages:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter,r *http.Request) {
		serveWs(hub,w,r)
	})

	s:= http.Server{Addr:":" + port, Handler:mux}
	return s.ListenAndServe()

}

func main(){
	port := flag.String("port", ":9000", "port for server")
	flag.Parse()
	runServer(*port)
}


