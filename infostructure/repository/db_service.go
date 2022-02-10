package repository

import (
	"github.com/yarikyarichek/streamer/entity"
	"gorm.io/gorm"
)

type Service interface {
	Migrate() error
	Create(msg *entity.Message) error
	Get(filter *entity.Message, offset, limit int) (*entity.Messages, error)
}

type dbService struct {
	*gorm.DB
}

func NewService(db *gorm.DB) Service { return &dbService{db} }

func (service *dbService) Migrate() error {
	if err := service.AutoMigrate(&entity.Message{}); err != nil {
		return err
	}
	return nil
}

func (service *dbService) Create(msg *entity.Message) error {
	return service.DB.Create(&msg).Error
}

func (service *dbService) Get(filter *entity.Message, offset, limit int) (*entity.Messages, error) {
	query := service.DB
	if offset > 0 {
		query.Offset(offset)
	}
	if limit > 0 {
		query.Limit(limit)
	}
	if filter != nil {
		query.Where(filter)
	}
	result := entity.Messages{}
	return &result, query.Find(&result).Error
}
