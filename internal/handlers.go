package internal

import (
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleSet(args []string) string {
	if len(args) < 2 {
		return "(error) ERR wrong number of arguments for SET command"
	}

	key := args[0]
	value := args[1]
	var expiration time.Time
	hasPX := false

	for i := 2; i < len(args); i++ {
		arg := strings.ToUpper(args[i])
		if arg == "PX" {
			if hasPX {
				return "(error) ERR PX already specified"
			}
			if i+1 >= len(args) {
				return "(error) ERR syntax error"
			}
			milisecondsStr := args[i+1]
			miliseconds, err := strconv.ParseInt(milisecondsStr, 10, 64)
			if err != nil || miliseconds < 0 {
				return "(error) ERR value is not an integer or out of range"
			}
			expiration = time.Now().Add(time.Duration(miliseconds) * time.Millisecond)
			i++
		} else if hasPX {
			return "(error) ERR syntax error after PX"
		} else {
			value += " " + args[i]
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = valueEntry{Value: value, Expiration: expiration}
	return "OK"
}

func (s *Server) handleGet(args []string) string {
	if len(args) != 1 {
		return "(error) ERR wrong number of arguments for GET command"
	}

	key := args[0]

	s.mu.RLock()
	entry, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return "(nil)"
	}

	if !entry.Expiration.IsZero() && time.Now().After(entry.Expiration) {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return "(nil)"
	}

	return entry.Value
}
