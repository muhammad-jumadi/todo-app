package datastore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	model "todo-app/model"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore() *DBStore {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		model.DBHost, model.DBPort, model.DBUser, model.DBPassword, model.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Successfully connected!")

	return &DBStore{
		db: db,
	}
}

func (ds *DBStore) GetCompleted(w http.ResponseWriter, r *http.Request) {
	var completed []model.TodoData

	query := `
		SELECT id, title, status
		FROM todo
		WHERE status = true
	`

	rows, err := ds.db.Query(query)
	if err != nil {
		log.Println("error on getting todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var data model.TodoData
		if err := rows.Scan(&data.ID, &data.Title, &data.Status); err != nil {
			log.Println("error on getting todo:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		completed = append(completed, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completed)
}

func (ds *DBStore) GetIncomplete(w http.ResponseWriter, r *http.Request) {
	var inCompleted []model.TodoData

	query := `
		SELECT id, title, status
		FROM todo
		WHERE status = false
	`

	rows, err := ds.db.Query(query)
	if err != nil {
		log.Println("error on getting todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var data model.TodoData
		if err := rows.Scan(&data.ID, &data.Title, &data.Status); err != nil {
			log.Println("error on getting todo:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		inCompleted = append(inCompleted, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inCompleted)
}

func (ds *DBStore) CreateTodo(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	query := `
		INSERT INTO todo (title) VALUES ($1)
	`

	stmt, err := ds.db.Prepare(query)
	if err != nil {
		log.Println("error on prepare statement todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer stmt.Close()

	_, err = stmt.Exec(title)
	if err != nil {
		log.Println("error on created todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (ds *DBStore) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	status, _ := strconv.ParseBool(r.FormValue("status"))
	query := `
		UPDATE todo SET status = $1
		WHERE id = $2
	`

	stmt, err := ds.db.Prepare(query)
	if err != nil {
		log.Println("error on prepare statement todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, id)
	if err != nil {
		log.Println("error on updated todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (ds *DBStore) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	query := `
		DELETE
		FROM todo
		WHERE id = $1
	`

	stmt, err := ds.db.Prepare(query)
	if err != nil {
		log.Println("error on prepare statement todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println("error on deleted todo:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
