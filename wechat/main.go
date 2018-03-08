package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	logLevel = "dev"
	port     = 8081
	token    = "kosfe9rtue9r1"
)


func main() {
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        &httpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}
	fmt.Printf("server addr:%s\n", server.Addr)
	log.Println(fmt.Sprintf("Listen: %d", port))
	log.Fatal(server.ListenAndServe())
}
