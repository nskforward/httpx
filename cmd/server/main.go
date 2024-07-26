package main

import (
	"fmt"

	"github.com/nskforward/httpx"
)

func main() {
	addr := ":80"
	var router httpx.Router
	router.Route("/", httpx.Echo)
	fmt.Println("ready to handle connections on", addr)
	router.Listen(addr)
}
