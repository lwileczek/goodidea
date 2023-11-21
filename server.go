package goodidea

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func listTasks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	var err error
	var taskList []Task
	if title == "" {
		taskList, err = getAllTasks()
	} else {
		taskList, err = getSomeTasks(title)
	}
	if err != nil {
		Logr.Error("Could not get tasks from db", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/list-tasks.html", "templates/task.html"))
	err = tmpl.Execute(w, taskList)
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

	score, err := updateTaskScore(uint32(id), r.FormValue(fmt.Sprintf("scorekeeper%d", id)) == "inc")
	if err != nil {
		fmt.Println("Could not update task score", err)
	}

	fmt.Fprintf(w, fmt.Sprintf("%d", score))
}

// createTask - recieve a title and body and create a new task in the DB with default values
// return an HTML block of the Task summarized to be displayed on landing page
func createTask(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) //10MB
	if err != nil {
		fmt.Println("Error parsing form", err)
	}

	title := r.FormValue("title")
	if title == "" {
		fmt.Fprintf(w, "<p>Tasks must have a title.</p>")
		return
	}

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

	//------------
	//TODO: Add filepath to the data record so when they open the task the image shows
	file, header, err := r.FormFile("taskImgs") //matches html
	if err != nil {
		Logr.Error("Could not retrive file from form data", "error", err)
		return
	}
	defer file.Close()
	nameComponents := strings.Split(header.Filename, ".")
	if len(nameComponents) != 2 {
		Logr.Error("Could not find the extension of the uploaded file", "error", err)
		return
	}

	b, err := io.ReadAll(file)
	if err != nil {
		Logr.Error("could not read bytes out of file sent", "error", err)
		return
	}
	m := NewFileManager()
	m.StoreFile(b, nameComponents[1])
	//------

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

	if r.FormValue("comments") == "" {
		return
	}
	var pu *string
	username := r.FormValue("username")
	if username != "" {
		pu = &username
	}

	if err := addComment(uint32(id), pu, r.FormValue("comments")); err != nil {
		Logr.Error("Could not save a new comment", "error", err)
		fmt.Fprintf(w, "<p>ERROR! could not insert the comment into the DB</p>")
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/make-comment.html", "templates/comment.html"))
	err = tmpl.Execute(w, Comment{TaskID: uint32(id), User: pu, Content: r.FormValue("comments"), CreatedAt: time.Now()})
	if err != nil {
		Logr.Error("could not render template for new comment", "error", err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("<h2>404 Could not find!</h2><p>Path Provided: %s</p>", r.URL)
	fmt.Fprintf(w, s)
}

func NewServer() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", index).Methods("GET")
	mux.HandleFunc("/prod", index).Methods("GET")
	mux.HandleFunc("/prod/goodidea", index).Methods("GET")
	mux.HandleFunc("/tasks", listTasks).Methods("GET")
	mux.HandleFunc("/tasks", createTask).Methods("POST")
	mux.HandleFunc("/tasks/{id}", viewTask).Methods("GET")
	mux.HandleFunc("/tasks/{id}/score", updateScore).Methods("POST")
	mux.HandleFunc("/tasks/{id}/comments", postComment).Methods("POST")

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	mux.PathPrefix("/static/").Handler(s)

	mux.NotFoundHandler = http.HandlerFunc(notFound)

	return mux
}
