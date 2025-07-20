package internal

import (
	"flag"
	"fmt"
	"os"
)

const DefaultPort = 8080

func FlagInit() int {
	help := flag.Bool("flag", false, "Show help message")
	port := flag.Int("port", DefaultPort, "Port number")

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	if *port < 0 || *port > 65535 {
		fmt.Fprintf(os.Stderr, "Error: port number %d is out of valid range\n", *port)
		os.Exit(1)
	}

	return *port
}
