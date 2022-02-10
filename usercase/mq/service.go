package mq

import (
	"io"
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
	cdone   chan int
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
		cdone:   make(chan int),
	}
}

func (s *service) Query() chan *entity.Message {
	return s.query
}

func (s *service) Size() int {
	return s.size
}
func (s *service) Clear() {}

func (s *service) Start() {}
