package main

import (
	"fmt"
	"net/http"
	"time"
)

func startClient(wakeup chan bool) {
	client := &http.Client{
		Timeout: time.Second * 10, // Устанавливаем таймаут в 10 секунд
	}
	for {
		funcClient(client)
		timesleep := nextTaskTime() - time.Now().UTC().Second()
		if timesleep > 0 {
			time.Sleep(time.Duration(timesleep) * time.Second)
		} else {
			<-wakeup
		}

	}
}
func stopClient(client *http.Client) {
	client.CloseIdleConnections()
}
func startServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/test", addHandler)
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
