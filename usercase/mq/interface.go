package mq

import "github.com/yarikyarichek/streamer/entity"

type Service interface {
	Query() chan *entity.Message
	Size() int
	Clear()
	Start()
}
