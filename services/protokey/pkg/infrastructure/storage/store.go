package storage

import (
	"protokey/pkg/app/service"
)

type Store struct {
	CommandChan  chan service.Command
	ResponseChan chan service.Response
}

func NewStore() *Store {
	commandChan := make(chan service.Command)
	responseChan := make(chan service.Response)

	s := &Store{
		CommandChan:  commandChan,
		ResponseChan: responseChan,
	}

	go s.run(commandChan, responseChan)

	return s
}

func (s *Store) run(commandChan <-chan service.Command, responseChan chan<- service.Response) {
	data := make(map[string]int)
	for command := range commandChan {
		switch command.Type {
		case service.SetOperation:
			data[command.Key] = command.Value
			responseChan <- service.Response{Err: nil}
		case service.GetOperation:
			val, ok := data[command.Key]
			if !ok {
				val = 0
			}
			responseChan <- service.Response{Value: val, Err: nil}
		case service.ListKeysOperation:
			var keys []string
			for k := range data {
				if len(k) >= len(command.Key) && k[:len(command.Key)] == command.Key {
					keys = append(keys, k)
				}
			}
			responseChan <- service.Response{Keys: keys, Err: nil}
		}
	}
}
