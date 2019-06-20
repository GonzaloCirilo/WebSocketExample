package main

import (
	"flag"
	"fmt"
	"net/http"
	"sockets/socket"
)

func runServer(port string) error{
	hub := socket.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter,r *http.Request) {
		socket.ServeWs(hub,w,r)
	})

	s:= http.Server{Addr:":" + port, Handler:mux}
	fmt.Printf("Server listening at port %s", port)
	return s.ListenAndServe()

}

func main(){
	port := flag.String("port", "9000", "port for server")
	flag.Parse()
	runServer(*port)
}


