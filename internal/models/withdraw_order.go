package models

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

type WithdrawOrder struct {
	Id                    uint            `gorm:"primary_key"`
	UserId                uint            `gorm:"column:userId"`
	OrderId               string          `gorm:"column:orderId"`
	OrderStatus           string          `gorm:"column:orderStatus"`
	ChannelCode           string          `gorm:"column:channelCode"`
	Amount                decimal.Decimal `gorm:"type:decimal(10,2)"`
	PoundageAmount        decimal.Decimal `gorm:"type:decimal(10,2);column:poundageAmount"`
	createTime            uint            `gorm:"column:createTime"`
	updateTime            uint64          `gorm:"column:updateTime"`
	ChannelFailRequestLog json.RawMessage `gorm:"column:channelFailRequestLog;type:json"`
	CreatedAt             time.Time       `gorm:"type:timestamp;column:createTime"`
	UpdatedAt             time.Time       `gorm:"type:timestamp;column:updateTime"`
}
