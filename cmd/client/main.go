package main

import (
	"fmt"

	"github.com/Ryanjones0911/chat/internal/network/client"
)

func main() {
	fmt.Println("Hello from client")
	client.Init()
}
