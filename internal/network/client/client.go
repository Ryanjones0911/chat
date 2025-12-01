// internal/network/client/client.go
package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Init connects to addr (e.g. "localhost:9000") and starts an
// interactive chat session with the server.
func Init(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not connect to %s: %v", addr, err)
	}
	fmt.Printf("Connected to %s\n", addr)
	fmt.Println("Type messages and press Enter. Ctrl+C to exit.")

	// Handle Ctrl+C / SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nClosing connection and exiting...")
		_ = conn.Close()
		os.Exit(0)
	}()

	// Goroutine: read from server and print to stdout
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("server read error: %v", err)
		}
		fmt.Println("Disconnected from server.")
		os.Exit(0)
	}()

	// Main goroutine: read from stdin and send to server
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		text := stdin.Text()
		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			log.Printf("write error: %v", err)
			break
		}
	}
	if err := stdin.Err(); err != nil {
		log.Printf("stdin error: %v", err)
	}

	_ = conn.Close()
}
