package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

func addHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	Priority := r.FormValue("Priority")
	Message := r.FormValue("Message")
	datetime := r.FormValue("datetime") //utc sec date

	if Priority == "" || Message == "" || datetime == "" {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	// Преобразование ID к целому числу
	PriorityInt, err := strconv.Atoi(Priority)
	if err != nil || PriorityInt > 3 {
		http.Error(w, "Invalid Priority", http.StatusBadRequest)
		return
	}
	datetimeInt, err := strconv.Atoi(datetime)
	datetimeInt /= 60
	if err != nil || datetimeInt < 0 {
		http.Error(w, "Invalid datetime", http.StatusBadRequest)
		return
	}

	go add_data(PriorityInt, Message, false, datetimeInt)
}
func add_data(PriorityInt int, Message interface{}, false bool, sendTime int) {
	new_msg := Messages{PriorityInt, Message, false}
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if PriorityInt == 3 {
		sm.priority_data = append(sm.priority_data, new_msg)
	} else {
		sm.data[sendTime] = append(sm.data[sendTime], new_msg)
	}

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	datetime := r.FormValue("datetime")

	datetimeInt, err := strconv.Atoi(datetime)

	if err != nil || datetimeInt < 0 {
		http.Error(w, "Invalid datetime", http.StatusBadRequest)
		return
	}
	datetimeInt /= 60
	send := send_data(datetimeInt)

	if len(send) == 0 {
		http.Error(w, "{}", http.StatusAccepted)
		return
	}
	for i, data := range send {
		if i != 0 {
			fmt.Println(data.Message, data.Priority)
		}

	}
	data := fmt.Sprintf("%v %v", send[0].Message, datetimeInt)
	fmt.Fprintf(w, data)
}

func filter(sm []Messages) (ret []Messages, ids []int) {
	for i, s := range sm {
		if !s.status {
			ret = append(ret, s)
			ids = append(ids, i)
		}

	}
	return
}

func send_data(send_date int) (sended []Messages) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	unsended_Priority, ids_Priority := filter(sm.priority_data)
	unsended_main, ids_main := filter(sm.data[send_date])

	if len(unsended_Priority) != 0 {
		sended = append(sended, unsended_Priority...)
		for _, val := range ids_Priority {
			sm.priority_data[val].status = true
		}
	}

	sort.Slice(unsended_main, func(i, j int) bool {
		return unsended_main[i].Priority > unsended_main[j].Priority
	})

	if len(unsended_main) != 0 {
		sended = append(sended, unsended_main...)
		for _, val := range ids_main {
			sm.data[send_date][val].status = true
		}
	}

	return
}
