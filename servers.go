package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type sendJson struct {
	message  interface{}
	priority int
	// datetime time.Time   `json:"time"`
}

func funcClient(client *http.Client) {

	datetimeInt := time.Now().UTC().Second() / 60
	send := send_data(datetimeInt)

	if len(send) == 0 {

		return
	}
	for i, data := range send {
		if i != 0 {
			fmt.Println(data.Message, data.Priority)
		}

	}
	send_data := sendJson{message: send[0].Message,
		priority: send[0].Priority,
		// datetime: send[0].datetime,
	}
	json_data, err := json.Marshal(send_data)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	url := "127.0.0.1:8088/task"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Status code:", resp.StatusCode)
	fmt.Println("Response body: ", string(body))
}

func startClient() {
	client := &http.Client{
		Timeout: time.Second * 10, // Устанавливаем таймаут в 10 секунд
	}
	for {
		funcClient(client)
		time.Sleep(1 * time.Second)
	}
}
func stopClient(client *http.Client) {
	client.CloseIdleConnections()
}
func startServer(port string) {
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
