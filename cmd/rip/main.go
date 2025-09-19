package main

import (
	"fmt"
	"rip/internal/api"
)

func main() {
	fmt.Println("Hello")

	api.StartServer()

	fmt.Println("Bye")
}
