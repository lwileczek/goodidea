package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {

	tsks, err := getAllTasks()
	if err != nil {
		fmt.Println("Could not get all tasks from db", err)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/task.html", "templates/footer.html", "templates/create.html"))

	err = tmpl.Execute(w, tsks)
	if err != nil {
		fmt.Println("Could not execute template", err)
	}

}

func updateScore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		fmt.Println("Could not parse id", err)
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form", err)
	}

	score, err := updateTaskScore(uint32(id), r.FormValue("scorekeeper") == "inc")
	if err != nil {
		fmt.Println("Could not update task score", err)
	}

	fmt.Fprintf(w, fmt.Sprintf("%d", score))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form", err)
	}

	title := r.FormValue("title")
	body := r.FormValue("details")
	err = addTask(title, body)
	if err != nil {
		fmt.Println("Could not update task score", err)
	}

	task := Task{
		Title: title,
		Body:  &body,
		Score: 0,
	}

	tmpl := template.Must(template.ParseFiles("templates/make-task.html", "templates/task.html"))
	err = tmpl.Execute(w, task)
	if err != nil {
		fmt.Println("Could not execute template", err)
	}
}

func NewServer() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/tasks/{id}/score", updateScore).Methods("POST")
	mux.HandleFunc("/tasks", createTask).Methods("POST")

	return mux
}
