package main

import (
	"fmt"
	"net/http"
	"time"
)

func startServerGet(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/get_data", getHandler)
	server := &http.Server{Addr: port, Handler: mux}
	fmt.Println("Starting server on port", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
func startServerPost(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addHandler)

	server := &http.Server{Addr: port, Handler: mux, WriteTimeout: 10 * time.Second}
	fmt.Println("Starting server on port", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func stopServer(server *http.Server) {
	fmt.Println("Stopping server...")
	err := server.Shutdown(nil)
	if err != nil {
		fmt.Println("Error stopping server:", err)
	}
	fmt.Println("Server stopped.")
}
