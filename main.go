package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type SyncData struct {
	data          map[int]([]Messages)
	priority_data []Messages
	mutex         sync.Mutex
}
type Messages struct {
	priority int
	message  interface{}
	status   bool
}

var sm = SyncData{data: make(map[int]([]Messages))}

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
