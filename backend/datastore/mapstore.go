package datastore

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Xanvial/todo-app-go/model"
	"github.com/gorilla/mux"
)

type MapStore struct {
	data map[string]bool
}

func NewMapStore() *MapStore {
	newData := make(map[string]bool, 0)

	return &MapStore{
		data: newData,
	}
}

func (ms *MapStore) GetCompleted(w http.ResponseWriter, r *http.Request) {
	// get completed data
	completed := make([]model.TodoData, 0)
	for k, v := range ms.data {
		if v {
			res := model.TodoData{
				Title:  k,
				Status: v,
			}

			completed = append(completed, res)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completed)
}

func (ms *MapStore) GetIncomplete(w http.ResponseWriter, r *http.Request) {
	// get incompleted data
	inCompleted := make([]model.TodoData, 0)
	for k, v := range ms.data {
		if !v {
			res := model.TodoData{
				Title:  k,
				Status: v,
			}

			inCompleted = append(inCompleted, res)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inCompleted)
}

func (ms *MapStore) CreateTodo(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	ms.data[title] = false
}

func (ms *MapStore) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	status, _ := strconv.ParseBool(r.FormValue("status"))

	ms.data[title] = status
}

func (ms *MapStore) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	delete(ms.data, title)
}
