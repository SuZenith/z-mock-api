package models

import "github.com/shopspring/decimal"

type RechargeOrders struct {
	Id          uint64          `gorm:"column:id;primary_key"`
	Amount      decimal.Decimal `gorm:"type:decimal(10,2)"`
	OrderNum    string          `gorm:"column:orderNum;type:varchar(255)"`
	OrderStatus string          `gorm:"column:orderStatus;type:varchar(255)"`
	UserId      uint            `gorm:"column:userId"`
	Mobile      string          `gorm:"column:mobile;type:varchar(255)"`
	PayCode     string          `gorm:"column:payCode;type:varchar(255)"`
	SubPayCode  string          `gorm:"column:subPayCode;type:varchar(255)"`
	UpdatedAt   uint            `gorm:"column:updatedAt"`
	CreatedAt   uint            `gorm:"column:createdAt"`
}
