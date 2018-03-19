package main

import (
	"log"
	"net"
	"os"
	"io"
)

func main() {
	done := make(chan bool)
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done<-true
	}()
	mustCopy(conn, os.Stdin)
	<-done
	conn.Close()
}


func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
