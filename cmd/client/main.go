// cmd/client/main.go
package main

import (
	"fmt"
	"os"

	"github.com/Ryanjones0911/chat/internal/network/client"
)

func main() {
	addr := "localhost:9000" // default
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	fmt.Println("Hello from client")
	client.Init(addr)
}
