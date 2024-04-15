package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func addHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	priority := r.FormValue("priority")
	message := r.FormValue("message")
	datetime := r.FormValue("datetime")

	if priority == "" || message == "" || datetime == "" {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	// Преобразование ID к целому числу
	priorityInt, err := strconv.Atoi(priority)
	if err != nil {
		http.Error(w, "Invalid priority", http.StatusBadRequest)
		return
	}
	datetimeInt, err := strconv.Atoi(datetime)
	if err != nil {
		http.Error(w, "Invalid datetime", http.StatusBadRequest)
		return
	}
	new_msg := Messages{priorityInt, message}
	go chan_add(new_msg, datetimeInt)
}
func chan_add(new_msg Messages, sendTime int) {

	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.data[sendTime] = append(sm.data[sendTime], new_msg)

}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	datetime := r.FormValue("datetime")

	datetimeInt, err := strconv.Atoi(datetime)
	if err != nil {
		http.Error(w, "Invalid priority", http.StatusBadRequest)
		return
	}
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	data := fmt.Sprintf("%v", sm.data[datetimeInt])
	fmt.Fprintf(w, data)
}
