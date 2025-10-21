package main

/*
* Server executable. No networking logic will be done
* here, this is just for the calling of functions and
* things like that
*
*
*
* */

import (
	"fmt"

	"chat/internal/network/server"
)

func main() {
	fmt.Println("Helllo from server")
	server.Init()
}
