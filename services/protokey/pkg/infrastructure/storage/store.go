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
	config            Config
	CommandChan       chan service.Command
	processedCommands []service.Command
	data              map[string]int
	dataFileHandle    *os.File
}

func NewStore(config Config) *Store {
	store := &Store{
		config:      config,
		CommandChan: make(chan service.Command),
		data:        make(map[string]int),
	}

	go store.initializeAndRun()
	return store
}

func (store *Store) initializeAndRun() {
	err := store.loadFromDisk()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: failed to load data file: %v\n", err)
	}

	store.openDataFileForAppend()
	defer store.closeDataFile()

	store.runLoop()
}

func (store *Store) loadFromDisk() error {
	filePath := filepath.Clean(store.config.DataFile)
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
			store.data[cmd.Key] = cmd.Value
		}
	}

	return scanner.Err()
}

func (store *Store) openDataFileForAppend() {
	f, err := os.OpenFile(store.config.DataFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)

	if os.IsNotExist(err) {
		f, err = os.Create(store.config.DataFile)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warning: cannot open data file for append: %v\n", err)
		return
	}

	store.dataFileHandle = f
}

func (store *Store) closeDataFile() {
	if store.dataFileHandle != nil {
		_ = store.dataFileHandle.Close()
	}
}

func (store *Store) runLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd, ok := <-store.CommandChan:
			if !ok {
				return
			}
			store.processCommand(cmd)

		case <-ticker.C:
			store.flushProcessedCommands()
		}
	}
}

func (store *Store) processCommand(cmd service.Command) {
	switch cmd.Type {
	case service.SetOperation:
		store.handleSetOperation(cmd)
	case service.GetOperation:
		store.handleGetOperation(cmd)
	case service.ListKeysOperation:
		store.handleListKeysOperation(cmd)
	default:
		cmd.Reply <- service.Response{Err: fmt.Errorf("unknown operation %d", cmd.Type)}
	}
}

func (store *Store) handleSetOperation(cmd service.Command) {
	store.data[cmd.Key] = cmd.Value
	store.processedCommands = append(store.processedCommands, cmd)
	cmd.Reply <- service.Response{Err: nil}
}

func (store *Store) handleGetOperation(cmd service.Command) {
	value, exists := store.data[cmd.Key]

	if !exists {
		value = 0
	}

	cmd.Reply <- service.Response{Value: value, Err: nil}
}

func (store *Store) handleListKeysOperation(cmd service.Command) {
	prefix := cmd.Key
	var keys []string

	for k := range store.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}

	cmd.Reply <- service.Response{Keys: keys, Err: nil}
}

func (store *Store) flushProcessedCommands() {
	if store.dataFileHandle == nil || len(store.processedCommands) == 0 {
		return
	}

	writer := bufio.NewWriter(store.dataFileHandle)
	for _, cmd := range store.processedCommands {
		line, err := json.Marshal(struct {
			Type  service.OperationType `json:"type"`
			Key   string                `json:"key"`
			Value int                   `json:"value"`
		}{
			Type:  cmd.Type,
			Key:   cmd.Key,
			Value: cmd.Value,
		})
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

	store.processedCommands = store.processedCommands[:0]
}
