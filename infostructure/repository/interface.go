package repository

import "github.com/yarikyarichek/streamer/entity"

type Service interface {
	Migrate() error
	Create(msgs *entity.Messages) error
	Get(filter *entity.Message, offset, limit int) (*entity.Messages, error)
}
