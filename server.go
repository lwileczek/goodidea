package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	tsks, err := getAllTasks()
	if err != nil {
		fmt.Println("Could not get all tasks from db", err)
	}

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	err = tmpl.ExecuteTemplate(w, "index.html", tsks)
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

//createTask - recieve a title and body and create a new task in the DB with default values
//return an HTML block of the Task summarized to be displayed on landing page
func createTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form", err)
	}

	title := r.FormValue("title")
	body := r.FormValue("details")
	taskID, err := addTask(title, body)
	if err != nil {
		Logr.Error("Could not create a new task", "error", err)
		fmt.Fprintf(w, "<p>Could not create a new task</p>")
		return
	}

	//On the summary page, don't show an entire task description which may be long
	if len(body) > 64 {
		body = body[:64]
	}

	task := Task{
		ID:    taskID,
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

func viewTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		fmt.Println("Could not parse id", err)
		fmt.Fprintf(w, "ERROR! Could not get task ID from request")
		return
	}
	tsk, err := getTasksByID(uint32(id))
	if err != nil {
		Logr.Error("could not get task by id", "taskID", id, "error", err)
		fmt.Fprintf(w, "Error could not get a task with this ID")
		return
	}

	comments, err := getAllTaskComments(uint32(id))
	if err != nil {
		Logr.Error("could not get comments for task", "taskID", id, "error", err)
	}

	tsk.Comments = comments

	tmpl := template.Must(template.ParseFiles("templates/show-task.html", "templates/header.html", "templates/meta.html", "templates/footer.html", "templates/comment.html"))
	err = tmpl.Execute(w, tsk)
	if err != nil {
		Logr.Error("could not render template for task", "taskID", id, "error", err)
	}
}

func postComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		Logr.Error("Could not parse id", "error", err)
		fmt.Fprintf(w, "ERROR! Could not get task ID from request")
		return
	}

	err = r.ParseForm()
	if err != nil {
		Logr.Error("Error parsing form", "error", err)
		fmt.Fprintf(w, "<p>ERROR! Could not pase the form data</p>")
		return
	}

	var username string = r.FormValue("username")
	//if r.FormValue("username") != "" {
	//	*username = r.FormValue("username")
	//}

	if err := addComment(uint32(id), &username, r.FormValue("comments")); err != nil {
		Logr.Error("Could not save a new comment", "error", err)
		fmt.Fprintf(w, "<p>ERROR! could not insert the comment into the DB</p>")
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/make-comment.html", "templates/comment.html"))
	err = tmpl.Execute(w, Comment{TaskID: uint32(id), User: &username, Content: r.FormValue("comments"), CreatedAt: time.Now()})
	if err != nil {
		Logr.Error("could not render template for new comment", "error", err)
	}
}

func NewServer() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/tasks/{id}", viewTask).Methods("GET")
	mux.HandleFunc("/tasks/{id}/score", updateScore).Methods("POST")
	mux.HandleFunc("/tasks/{id}/comments", postComment).Methods("POST")
	mux.HandleFunc("/tasks", createTask).Methods("POST")

	return mux
}
