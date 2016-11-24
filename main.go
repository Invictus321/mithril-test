package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Task struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

var tasks []Task

func setDoneById(id int64, done bool) error {
	for _, task := range tasks {
		if id == task.Id {
			task.Done = done
			return nil
		}
	}
	return fmt.Errorf("Task not found")
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST")
	switch r.Method {
	case "GET":
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
	case "PUT":
		id := r.FormValue("id")
		done := r.FormValue("done")
		if id == "" || done == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			doneBool, err := strconv.ParseBool(done)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if err := setDoneById(idInt, doneBool); err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	case "POST":
		id := r.FormValue("id")
		description := r.FormValue("description")
		if description == "" || id == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			idInt, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			task := Task{
				idInt,
				description,
				false,
			}
			tasks = append(tasks, task)
		}
	case "OPTIONS":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	tasks = []Task{
		Task{0, "Do stuff", false},
		Task{1, "Be awesome", true},
	}
	http.HandleFunc("/tasks", handle)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":3000", nil)
}
