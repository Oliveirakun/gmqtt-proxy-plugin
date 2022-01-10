package proxy

import "sync"

type Semaphore struct {
	locker *sync.RWMutex
	open   bool
}

func NewSemaphore() *Semaphore {
	return &Semaphore{
		locker: &sync.RWMutex{},
		open:   true,
	}
}

func (s *Semaphore) IsOpen() bool {
	s.locker.RLock()
	defer s.locker.RUnlock()

	return s.open
}

func (s *Semaphore) Open() {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.open = true
}

func (s *Semaphore) Close() {
	s.locker.Lock()
	defer s.locker.Unlock()

	s.open = false
}
