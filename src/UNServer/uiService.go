package main

import (
	"html/template"
	"log"
	"net/http"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.New("index.tmpl").Delims("[[", "]]").ParseFiles("./templates/index.tmpl")
	tmpl, err := template.New("index.html").ParseGlob("./templates/*")
	if err != nil {
		log.Println(err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
