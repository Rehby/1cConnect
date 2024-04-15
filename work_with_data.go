package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func addHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	Priority := r.FormValue("priority")
	Message := r.FormValue("message")
	datetime := r.FormValue("datetime") //unix date

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
	datetimeInt, err := strconv.ParseInt(datetime, 10, 64)
	datetimeInt /= 60
	if err != nil || datetimeInt < 0 {
		http.Error(w, "Invalid datetime", http.StatusBadRequest)
		return
	}

	go add_data(PriorityInt, Message, false, datetimeInt)
}
func add_data(PriorityInt int, Message interface{}, false bool, sendTime int64) {
	new_msg := Messages{PriorityInt, Message, false}
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if PriorityInt == 3 {
		sm.priority_data = append(sm.priority_data, new_msg)
	} else {
		sm.data[sendTime] = append(sm.data[sendTime], new_msg)
	}

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

func send_data(send_date int64) (sended []Messages) {
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

func nextTaskTime() int64 {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	dataprior, _ := filter(sm.priority_data)
	if len(dataprior) != 0 {
		return 0
	}

	if len(sm.data) == 0 {
		return 0
	}

	// i := 0

	var keys []int64
	for key, _ := range sm.data {
		fmt.Println(time.Now().Unix() / 60)
		if key <= time.Now().Unix()/60 {
			keys = append(keys, key)
		}

	}
	fmt.Println(keys[0]*60 - time.Now().Unix())
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys[0]*60 - time.Now().Unix()
}

func funcClient(client *http.Client) {

	datetimeInt := time.Now().Unix() / 60

	send := send_data(datetimeInt)

	if len(send) == 0 {

		return
	}
	for i, data := range send {
		if i != 0 {
			fmt.Println(data.Message, data.Priority)
		}

	}

	json_data, err := json.Marshal(send[0])
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	url := "http://127.0.0.1:8080/task"
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
