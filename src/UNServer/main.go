package main

import (
	ctx "context"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"
	// "code.google.com/p/gcfg"
)

func main() {
	ctx.InitServerContext()

	runtime.GOMAXPROCS(12)
	port := ":8090"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Init Listen error: ", err)
		time.Sleep(1)
		return
	}

	http.HandleFunc("/index", indexHandler)
	log.Printf("UN HTTP Server is listenting on %s\n", port)
	http.Serve(lis, nil)
}
