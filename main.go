package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	switch r.Method {
	case "GET":
		//create some fake tasks
		tasks := []Task{
			Task{"Do stuff", false},
			Task{"Be awesome", true},
		}
		response := struct {
			Tasks []Task `json:"tasks"`
		}{
			tasks,
		}
		b, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(b)
		}
	case "OPTIONS":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/tasks", handle)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}
