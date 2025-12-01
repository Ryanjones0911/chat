package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

// we need a startup function that has a function listening on some import
func StartServer(portNum string) {
	// seed RNG once
	rand.Seed(time.Now().UnixNano())

	listen, err := net.Listen("tcp", portNum)

	// this is to listen for system interrupts and trigger clean shutdown upon detection
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, os.Interrupt)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Now listening on port %s\nPress Ctrl + C to safely shutdown\n", portNum)
	}

	// shared list of clients + names + colors
	var (
		mu      sync.Mutex
		clients = make(map[net.Conn]struct{})
		names   = make(map[net.Conn]string)
		colors  = make(map[net.Conn]string)
	)

	// color choices and reset
	colorChoices := []string{
		"\033[31m", // red
		"\033[32m", // green
		"\033[33m", // yellow
		"\033[34m", // blue
		"\033[35m", // magenta
		"\033[36m", // cyan
		"\033[91m", // bright red
		"\033[92m", // bright green
		"\033[93m", // bright yellow
		"\033[94m", // bright blue
		"\033[95m", // bright magenta
		"\033[96m", // bright cyan
	}
	reset := "\033[0m"

	go func() {
		<-shutdownSig
		fmt.Println("\nServer shutting down...")
		listen.Close()

		mu.Lock()
		for c := range clients {
			c.Close()
		}
		mu.Unlock()

		os.Exit(0)
	}()

	for {
		connection, err := listen.Accept()
		if err != nil {
			return // don't fatal; accept fails normally on shutdown
		}

		mu.Lock()
		clients[connection] = struct{}{}
		mu.Unlock()

		go func(conn net.Conn) {
			// scanner for this connection
			scanner := bufio.NewScanner(conn)

			// ask for name
			io.WriteString(conn, "Enter your name: ")
			if !scanner.Scan() {
				conn.Close()
				return
			}
			name := scanner.Text()

			// assign random color
			color := colorChoices[rand.Intn(len(colorChoices))]

			mu.Lock()
			names[conn] = name
			colors[conn] = color
			mu.Unlock()

			// announce join
			joinMsg := fmt.Sprintf("*** %s%s%s joined the chat ***\n", color, name, reset)
			mu.Lock()
			for c := range clients {
				io.WriteString(c, joinMsg)
			}
			mu.Unlock()

			defer func() {
				// grab values before removing
				mu.Lock()
				leftName := names[conn]
				leftColor := colors[conn]
				delete(names, conn)
				delete(colors, conn)
				delete(clients, conn)
				mu.Unlock()
				conn.Close()

				// announce leave to remaining clients
				leaveMsg := fmt.Sprintf("*** %s%s%s left the chat ***\n", leftColor, leftName, reset)
				mu.Lock()
				for c := range clients {
					io.WriteString(c, leaveMsg)
				}
				mu.Unlock()
			}()

			// main chat loop
			for scanner.Scan() {
				text := scanner.Text()
				msg := fmt.Sprintf("%s%s%s: %s\n", colors[conn], names[conn], reset, text)

				mu.Lock()
				for c := range clients {

					io.WriteString(c, msg)
				}
				mu.Unlock()
			}
		}(connection)
	}
}
