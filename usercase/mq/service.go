package mq

import (
	"io"
	"log"
	"runtime"
	"sync"

	"github.com/yarikyarichek/streamer/entity"
	"github.com/yarikyarichek/streamer/infostructure/repository"
)

type service struct {
	repo    repository.Service
	size    int
	maxProc int
	tmp     int
	mux     sync.Mutex
	w       io.Writer
	storage []entity.Messages
	query   chan *entity.Message
	done    chan int
}

func NewService(size, maxProc int, w io.Writer, repo repository.Service) Service {

	if size < 1 {
		size = 1
	}
	maxProcRuntime := runtime.GOMAXPROCS(maxProc)
	if maxProc < 1 {
		maxProc = maxProcRuntime
	}

	storage := make([]entity.Messages, maxProc)
	for i := range storage {
		storage[i] = make(entity.Messages, 0, size)
	}

	return &service{
		repo:    repo,
		size:    size,
		maxProc: maxProc,
		tmp:     0,
		mux:     sync.Mutex{},
		w:       w,
		storage: storage,
		query:   make(chan *entity.Message, 100000),
		done:    make(chan int),
	}
}

func (s *service) Query() chan *entity.Message {
	return s.query
}

func (s *service) Size() int {
	return s.size
}

func (s *service) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	for i := range s.storage {
		err := s.repo.Create(&s.storage[i])
		if err != nil {
			log.Println(err)
		}
		s.storage[i] = make(entity.Messages, 0, s.size)
	}
}

func (s *service) Start() {
	go s.schedule()
	for i := range s.storage {
		s.done <- i
	}
}

func (s *service) schedule() {
	for {
		i, opened := <-s.done
		if !opened {
			return
		}
		go s.run(i)
	}
}

func (s *service) run(i int) {
	for {
		s.storage[i] = append(s.storage[i], <-s.query)
		if len(s.storage[i]) == s.size {
			break
		}
	}
	s.mux.Lock()
	err := s.repo.Create(&s.storage[i])
	if err != nil {
		log.Println(err)
	}
	s.storage[i] = make(entity.Messages, 0, s.size)
	s.mux.Unlock()
	s.done <- i
}
