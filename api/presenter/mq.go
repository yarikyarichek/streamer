package presenter

import (
	"time"

	"github.com/google/uuid"
	"github.com/yarikyarichek/streamer/entity"
)

const MinLimitSize = 1
const MaxLimitSize = 10000

type GetMessageRequest struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Tag       string    `json:"tag"`
	Text      string    `json:"text"`
	Offset    int       `json:"offset"`
	Limit     int       `json:"limit"`
}

type GetMessageResponse entity.Messages

type CreateMessegeRequest struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type CreateMessegeRequests []CreateMessegeRequest

type CreateMessegeResponse struct {
	Status string `json:"status"`
}

func (mr *GetMessageRequest) ToMessage() *entity.Message {
	return &entity.Message{
		ID:        mr.ID,
		CreatedAt: mr.CreatedAt,
		Tag:       mr.Tag,
		Text:      mr.Text,
	}
}

func (mr *GetMessageRequest) ValidateLimit() int {
	if mr.Limit < MinLimitSize && mr.Limit > MaxLimitSize {
		return MaxLimitSize
	}
	return mr.Limit
}

func (mr *CreateMessegeRequest) ToMessage() *entity.Message {
	return &entity.Message{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Tag:       mr.Tag,
		Text:      mr.Text,
	}
}

func (mr *CreateMessegeRequests) ToMessage() *entity.Messages {
	var result entity.Messages = make(entity.Messages, 0, len(*mr))
	for _, m := range *mr {
		result = append(result, m.ToMessage())
	}
	return &result
}
