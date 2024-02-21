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

var tmpl *template.Template

func index(w http.ResponseWriter, r *http.Request) {
	tsks, err := getAllTasks(25)
	if err != nil {
		Logr.Error("Could not get all tasks from db", err)
	}

	err = tmpl.ExecuteTemplate(w, "index.html", tsks)
	if err != nil {
		Logr.Error("Could not execute template", err)
	}
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	var err error
	var taskList []Task
	if title == "" {
		taskList, err = getAllTasks(35)
	} else {
		taskList, err = getSomeTasks(title)
	}
	if err != nil {
		Logr.Error("Could not get tasks from db", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "list-tasks.html", taskList)
	if err != nil {
		Logr.Error("Could not execute template", err)
	}
}

func updateScore(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		Logr.Error("Could not parse id", err)
	}

	err = r.ParseForm()
	if err != nil {
		Logr.Error("Error parsing form", err)
	}

	score, err := updateTaskScore(uint32(id), r.FormValue(fmt.Sprintf("scorekeeper%d", id)) == "inc")
	if err != nil {
		Logr.Error("Could not update task score", err)
	}

	fmt.Fprintf(w, fmt.Sprintf("%d", score))
}

// createTask - recieve a title and body and create a new task in the DB with default values
// return an HTML block of the Task summarized to be displayed on landing page
func createTask(w http.ResponseWriter, r *http.Request) {
	//TODO: MultipartReader to transform this to a steam
	err := r.ParseMultipartForm(32 << 20) //32MB
	if err != nil {
		Logr.Error("Error parsing form", err)
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

	err = tmpl.ExecuteTemplate(w, "make-task.html", task)
	if err != nil {
		Logr.Error("Could not execute template", err)
	}

	//Add any images which may have been sent along with the form data
	fhs, ok := r.MultipartForm.File["taskImgs"] //matches html
	if !ok {
		return
	}

	filePaths := make([]string, len(fhs))
	for i, fh := range fhs {
		f, err := fh.Open()
		defer f.Close()

		nameComponents := strings.Split(fh.Filename, ".")
		if len(nameComponents) != 2 {
			Logr.Error("Could not find the extension of the uploaded file", "error", err)
			return
		}

		b, err := io.ReadAll(f)
		if err != nil {
			Logr.Error("could not read bytes out of file sent", "error", err)
			return
		}

		//TODO: can the rest of the following be done in a go routine?
		m := NewFileManager()
		s, err := m.StoreFile(b, nameComponents[1])
		if err != nil {
			Logr.Error("Unable to store images", "task", taskID, "error", err.Error())
			return
		}

		f.Close()
		filePaths[i] = s
	}
	go saveTaskImages(taskID, filePaths)
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		Logr.Error("Could not parse id", err)
		fmt.Fprintf(w, "ERROR! Could not get task ID from request")
		return
	}
	//TODO: Why not get all of the comments for the task in one query?
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

	err = tmpl.ExecuteTemplate(w, "show-task.html", tsk)
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

	err = tmpl.ExecuteTemplate(
		w,
		"make-comment.html",
		Comment{TaskID: uint32(id), User: pu, Content: r.FormValue("comments"), CreatedAt: time.Now()},
	)
	if err != nil {
		Logr.Error("could not render template for new comment", "error", err)
	}
}

func displayTaskImages(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		Logr.Error("Could not parse id", "err", err)
		fmt.Fprintf(w, "ERROR! Could not get task ID from request")
		return
	}
	paths, err := getTaskImages(uint32(id))
	if err != nil {
		Logr.Error("could not get image paths", "task", id, "err", err)
		fmt.Fprintf(w, "ERROR! Could not get images for task %d", id)
		return
	}
	if len(paths) == 0 {
		return
	}
	//TODO: Move to a template so tailwind will find the css classes
	content := ""
	for _, p := range paths {
		content += fmt.Sprintf(`<img onclick="enlargeModal()" class="h-32 w-32 mx-5 border-2 border-sky-900 cursor-pointer" src="%s" alt="task-image" width="128" height="128">`, p)
	}
	//This script is defiend in src/showTask.js, it adds a listener to each image
	content += `<script src="/static/enlargeImages.js"></script>`
	fmt.Fprintf(w, content)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("<h2>404 Could not find!</h2><p>Path Provided: %s</p>", r.URL)
	fmt.Fprintf(w, s)
}

func NewServer() *mux.Router {
	//Setup the templates so the endpoints work
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	mux := mux.NewRouter()
	mux.HandleFunc("/", index).Methods("GET")
	mux.HandleFunc("/prod", index).Methods("GET")
	mux.HandleFunc("/prod/goodidea", index).Methods("GET")
	mux.HandleFunc("/tasks", listTasks).Methods("GET")
	mux.HandleFunc("/tasks", createTask).Methods("POST")
	mux.HandleFunc("/tasks/{id}", viewTask).Methods("GET")
	mux.HandleFunc("/tasks/{id}/score", updateScore).Methods("POST")
	mux.HandleFunc("/tasks/{id}/comments", postComment).Methods("POST")
	mux.HandleFunc("/tasks/{id}/images", displayTaskImages).Methods("GET")

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	mux.PathPrefix("/static/").Handler(s)

	mux.NotFoundHandler = http.HandlerFunc(notFound)

	return mux
}
