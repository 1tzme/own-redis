package internal

import "fmt"

func usage() {
	fmt.Println(`Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help

Options:
  --help       Show this screen.
  --port N     Port number.`)
}
