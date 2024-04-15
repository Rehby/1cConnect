package main

import (
	"fmt"
	"net/http"
	"sync"
)

type SyncData struct {
	data          map[int]([]Messages)
	priority_data []Messages
	mutex         sync.Mutex
}
type Messages struct {
	Priority int
	Message  interface{}
	status   bool
}

var sm = SyncData{data: make(map[int]([]Messages))}

// var queue = make(map[int](chan Messages))

func main() {
	portPost := ":8080"

	// Запуск сервера в отдельной горутине
	go func() {
		startServer(portPost)
	}()

	go func() {

		startClient()
	}()

	// Отключение сервера при завершении main
	defer stopServer(&http.Server{Addr: portPost})

	// Ожидание нажатия Enter для завершения программы
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
