//go:build ignore
package main

import (
	"errors"
	"html/template"
	"log"
	"os"
)

func makeLogin() error {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	f, err := os.Create("./pages/login.html")
	if err != nil {
		log.Println("create file error", err)
		return err
	}
	defer f.Close()
	tmpl.ExecuteTemplate(f, "login.html", nil)

	return nil
}

func createDir(p string) error {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(p, os.ModePerm)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {
	if err := createDir("pages"); err != nil {
		log.Println("Couldn't create a directory `pages`")
		panic(err)
	}
	if err := makeLogin(); err != nil {
		log.Println("couldn't create the login page!")
		panic(err)
	}
}
