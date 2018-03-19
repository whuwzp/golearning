package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	var wg sync.WaitGroup

	input := bufio.NewScanner(c)
	for input.Scan() {
		go func(c net.Conn, shout string, delay time.Duration) {
			fmt.Fprintln(c, "\t", strings.ToUpper(shout))
			time.Sleep(delay)
			fmt.Fprintln(c, "\t", shout)
			time.Sleep(delay)
			fmt.Fprintln(c, "\t", strings.ToLower(shout))

		}(c, input.Text(), 1*time.Second)
		wg.Add(1)
	}
	// NOTE: ignoring potential errors from input.Err()
	wg.Wait()
	c.Close()
}

