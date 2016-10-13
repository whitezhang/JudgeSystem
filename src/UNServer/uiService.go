package main

import (
	//"compress/gzip"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.New("index.tmpl").Delims("[[", "]]").ParseFiles("./templates/index.tmpl")
	/*
		w.Header().Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
	*/

	tmpl, err := template.New("index.html").ParseGlob("./templates/*")
	if err != nil {
		log.Println(err)
		return
	}
	err = tmpl.Execute(w, nil)
	//err = tmpl.Execute(gz, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
