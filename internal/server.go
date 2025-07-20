package internal

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	port int
	data map[string]valueEntry
	mu   sync.RWMutex
}

type valueEntry struct {
	Value      string
	Expiration time.Time
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
		data: make(map[string]valueEntry),
	}
}

func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen UDP: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 4096)
	fmt.Printf("Server listening on %s\n", conn.LocalAddr().String())

	s.expiredKeysCleanup()

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		go s.handleRequest(conn, clientAddr, buffer[:n])
	}
}

func (s *Server) handleRequest(conn *net.UDPConn, clientAddr *net.UDPAddr, message []byte) {
	request := strings.TrimSpace(string(message))
	parts := strings.Fields(request)
	response := ""

	if len(parts) == 0 {
		response = "(error) ERR unknown command"
	} else {
		command := strings.ToUpper(parts[0])

		switch command {
		case "PING":
			response = "PONG"
		case "SET":
			response = s.handleSet(parts[1:])
		case "GET":
			response = s.handleGet(parts[1:])
		default:
			response = fmt.Sprintf("(error) ERR unknown command %s", parts[0])
		}
	}

	_, err := conn.WriteToUDP([]byte(response), clientAddr)
	if err != nil {
		fmt.Printf("Error sending response to %v: %v\n", clientAddr, err)
	}
}

func (s *Server) expiredKeysCleanup() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			now := time.Now()
			s.mu.Lock()
			for key, entry := range s.data {
				if !entry.Expiration.IsZero() && now.After(entry.Expiration) {
					delete(s.data, key)
				}
			}
			s.mu.Unlock()
		}
	}()
}
