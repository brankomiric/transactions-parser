package parser

import (
	"sync"
)

// Just a simple in-memory storage
type storage struct {
	lock    sync.Mutex
	records map[string][]Transaction
}

func (s *storage) Add(address string, tx Transaction) {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, exists := s.records[address]
	if exists {
		s.records[address] = append(s.records[address], tx)
	} else {
		s.records[address] = []Transaction{tx}
	}
}

func (s *storage) Get(address string) []Transaction {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.records[address]
}
