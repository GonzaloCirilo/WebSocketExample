package main

import (
	"flag"
	"net/http"
	"sockets/socket"
)

func runServer(port string) error{
	hub := socket.NewHub()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter,r *http.Request) {
		socket.ServeWs(hub,w,r)
	})

	s:= http.Server{Addr:":" + port, Handler:mux}
	return s.ListenAndServe()

}

func main(){
	port := flag.String("port", "9000", "port for server")
	flag.Parse()
	runServer(*port)
}


