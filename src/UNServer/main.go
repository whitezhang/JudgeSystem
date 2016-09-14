package main

import (
	ctx "context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	// "code.google.com/p/gcfg"
)

func main() {
	ctx.InitServerContext()
	initSessionManager()

	runtime.GOMAXPROCS(ctx.SvrCtx.SvrCfg.Server.NumCPU)
	port := fmt.Sprintf(":%d", ctx.SvrCtx.SvrCfg.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Init Listen error: ", err)
		return
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexPageHandler)
	/*
		http.HandleFunc("/index", indexPageHandler)
		http.HandleFunc("/contests", contestPageHandler)
		http.HandleFunc("/problems", problemPageHandler)
		http.HandleFunc("/status", statusPageHandler)
		http.HandleFunc("/ratings", ratingPageHandler)
		http.HandleFunc("/about", aboutPageHandler)
		http.HandleFunc("/submit", submitPageHandler)
		http.HandleFunc("/problem", singleProblemHandler)
		http.HandleFunc("/contestinfo", singleContestHandler)
	*/

	http.HandleFunc("/sislogin", sessionHandler)
	http.HandleFunc("/slogin", loginHandler)
	http.HandleFunc("/slogout", logoutHandler)
	http.HandleFunc("/sregister", registHandler)
	http.HandleFunc("/suser", userHandler)
	http.HandleFunc("/sstatusinfo", statusInfoHandler)
	http.HandleFunc("/sprobleminfo", problemInfoHandler)
	http.HandleFunc("/sproblems", problemsHandler)
	http.HandleFunc("/scontestinfo", contestInfoHandler)
	http.HandleFunc("/scontests", contestsHandler)

	http.HandleFunc("/ssubmit", submitHandler)

	log.Printf("UN HTTP Server is listenting on %s\n", port)
	http.Serve(lis, nil)
}
