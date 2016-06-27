package main

import (
	ctx "context"
	"log"
	"net"
	"net/http"
	"runtime"
	// "code.google.com/p/gcfg"
)

func main() {
	ctx.InitServerContext()

	runtime.GOMAXPROCS(12)
	port := ":8090"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Init Listen error: ", err)
		return
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/index", indexPageHandler)
	http.HandleFunc("/contests", contestPageHandler)
	http.HandleFunc("/problems", problemPageHandler)
	http.HandleFunc("/status", statusPageHandler)
	http.HandleFunc("/ratings", ratingPageHandler)
	http.HandleFunc("/about", aboutPageHandler)
	http.HandleFunc("/submit", submitPageHandler)

	http.HandleFunc("/problem", singleProblemHandler)
	http.HandleFunc("/contestinfo", singleContestHandler)

	http.HandleFunc("/slogin", loginHandler)
	http.HandleFunc("/suser", userHandler)
	http.HandleFunc("/sproblem", problemHandler)
	http.HandleFunc("/scontestinfo", contestInfoHandler)
	http.HandleFunc("/scontests", contestsHandler)
	log.Printf("UN HTTP Server is listenting on %s\n", port)
	http.Serve(lis, nil)
}
