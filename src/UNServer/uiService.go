package main

import (
	"html/template"
	"log"
	"net/http"
)

func singleContestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("contest.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func singleProblemHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("problem.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func statusPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("status.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func ratingPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("ratings.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func problemPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("problems.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func aboutPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("about.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func contestPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("contests.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.New("index.tmpl").Delims("[[", "]]").ParseFiles("./templates/index.tmpl")
	tmpl, err := template.New("index.tmpl").Delims("[[", "]]").ParseGlob("./templates/*")
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
