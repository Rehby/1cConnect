package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type SyncMap struct {
	data  map[int]([]Messages)
	mutex sync.Mutex
}
type Messages struct {
	priority int
	message  interface{}
}

var sm = SyncMap{data: make(map[int]([]Messages))}

// var queue = make(map[int](chan Messages))

func read_data(val int) (err error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	value, ok := sm.data[val]
	if !ok {
		fmt.Printf("not found")
	} else {

		fmt.Printf(strconv.Itoa(value[0].priority)+" %v", value[0].message)
	}
	return err
}

func startServerGet(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", testHandler)
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

func main() {
	portPost := ":8080"
	portGet := ":8081"
	// Запуск сервера в отдельной горутине
	go func() {
		startServerPost(portPost)
	}()
	go func() {
		startServerGet(portGet)
	}()

	// Отключение сервера при завершении main
	defer stopServer(&http.Server{Addr: portPost})
	defer stopServer(&http.Server{Addr: portGet})
	// Ожидание нажатия Enter для завершения программы
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
