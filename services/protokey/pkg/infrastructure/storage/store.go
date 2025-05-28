package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"protokey/pkg/app/service"
)

type Store struct {
	config          Config
	CommandChan     chan service.Command
	ResponseChan    chan service.Response
	pendingCommands []service.Command
	data            map[string]int
	dataFileHandle  *os.File
}

func NewStore(config Config) *Store {
	store := &Store{
		config:       config,
		CommandChan:  make(chan service.Command),
		ResponseChan: make(chan service.Response),
		data:         make(map[string]int),
	}

	go store.initializeAndRun()
	return store
}

func (s *Store) initializeAndRun() {
	if err := s.loadFromDisk(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: failed to load data file: %v\n", err)
	}

	s.openDataFileForAppend()
	defer s.closeDataFile()

	s.runLoop()
}

func (s *Store) loadFromDisk() error {
	filePath := filepath.Clean(s.config.DataFile)
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var cmd service.Command
		if err := json.Unmarshal(scanner.Bytes(), &cmd); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: invalid command in data file: %v\n", err)
			continue
		}
		if cmd.Type == service.SetOperation {
			s.data[cmd.Key] = cmd.Value
		}
	}

	return scanner.Err()
}

func (s *Store) openDataFileForAppend() {
	f, err := os.OpenFile(s.config.DataFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)

	if os.IsNotExist(err) {
		f, err = os.Create(s.config.DataFile)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: cannot open data file for append: %v\n", err)
		return
	}

	s.dataFileHandle = f
}

func (s *Store) closeDataFile() {
	if s.dataFileHandle != nil {
		_ = s.dataFileHandle.Close()
	}
}

func (s *Store) runLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd, ok := <-s.CommandChan:
			if !ok {
				return
			}
			s.processCommand(cmd)

		case <-ticker.C:
			s.flushPendingCommands()
		}
	}
}

func (s *Store) processCommand(cmd service.Command) {
	switch cmd.Type {
	case service.SetOperation:
		s.handleSetOperation(cmd)
	case service.GetOperation:
		s.handleGetOperation(cmd)
	case service.ListKeysOperation:
		s.handleListKeysOperation(cmd)
	default:
		s.ResponseChan <- service.Response{Err: fmt.Errorf("unknown operation %d", cmd.Type)}
	}
}

func (s *Store) handleSetOperation(cmd service.Command) {
	s.data[cmd.Key] = cmd.Value
	s.pendingCommands = append(s.pendingCommands, cmd)
	s.ResponseChan <- service.Response{Err: nil}
}

func (s *Store) handleGetOperation(cmd service.Command) {
	value, exists := s.data[cmd.Key]

	if !exists {
		value = 0
	}

	s.ResponseChan <- service.Response{Value: value, Err: nil}
}

func (s *Store) handleListKeysOperation(cmd service.Command) {
	var keys []string

	prefix := cmd.Key
	for k := range s.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}

	s.ResponseChan <- service.Response{Keys: keys, Err: nil}
}

func (s *Store) flushPendingCommands() {
	if s.dataFileHandle == nil || len(s.pendingCommands) == 0 {
		return
	}

	writer := bufio.NewWriter(s.dataFileHandle)
	for _, cmd := range s.pendingCommands {
		line, err := json.Marshal(cmd)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshaling command for persistence: %v\n", err)
			continue
		}

		_, err = writer.Write(line)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error writing command: %v\n", err)
			return
		}

		err = writer.WriteByte('\n')
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error writing command: %v\n", err)
			return
		}
	}

	if err := writer.Flush(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error flushing data file buffer: %v\n", err)
	}

	s.pendingCommands = s.pendingCommands[:0]
}
