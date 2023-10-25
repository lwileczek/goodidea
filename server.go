package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {

	//todos, err := model.GetAllTodos()
	//if err != nil {
	//	fmt.Println("Could not get all todos from db", err)
	//}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Could not execute template", err)
	}

}

func NewServer() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", index)

	return mux
}
