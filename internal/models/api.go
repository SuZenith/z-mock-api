package models

import (
	"encoding/json"
	"time"
)

type Api struct {
	Id           uint64          `gorm:"column:id;primary_key;"`
	UserId       string          `gorm:"column:user_id;not null;"`
	Uuid         string          `gorm:"column:uuid;not null;type:varchar(255)"`
	Path         string          `gorm:"column:path;not null;type:varchar(1024)"`
	Method       string          `gorm:"column:method;not null;type:varchar(32)"`
	StatusCode   int16           `gorm:"column:status_code;not null;type:int"`
	ContentType  string          `gorm:"column:content_type;not null;type:varchar(128)"`
	Headers      json.RawMessage `gorm:"column:headers;type:json"`
	ResponseBody string          `gorm:"column:response_body;not null;type:text"`
	CreatedAt    time.Time       `gorm:"column:created_at;not null;type:timestamp"`
	UpdatedAt    time.Time       `gorm:"column:updated_at;not null;type:timestamp"`
}
