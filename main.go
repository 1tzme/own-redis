package main

import (
	"fmt"
	"os"

	"own-redis/internal"
)

func main() {
	port := internal.FlagInit()

	server := internal.NewServer(port)
	err := server.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
