package main

import (
	"fmt"

	"chat/internal/network/client"
)

func main() {
	fmt.Println("Hello client")
	client.Init()
}
