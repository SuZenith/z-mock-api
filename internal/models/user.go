package models

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Users struct {
	Id                            uint `gorm:"primary_key"`
	Mobile                        string
	UserId                        uint `gorm:"column:userId"`
	Name                          string
	Avatar                        string
	ActivityAvatar                json.RawMessage `gorm:"type:json"`
	Password                      string
	FundPassword                  string
	VersionCode                   string
	Channel                       string
	Ip                            string
	NoRealNamePayCode             string
	Adcode                        int
	Nation                        string
	Province                      string `gorm:"column:provice"`
	City                          string
	District                      string
	MobileArea                    string
	CreatedAt                     uint
	UpdatedAt                     uint
	Imei                          string
	Oaid                          string
	Status                        string
	FeeContribute                 string
	Label                         json.RawMessage `gorm:"type:json"`
	ExperienceLevel               string
	Dan                           int
	DanExpireDate                 string
	NotAllowFollowedUpFlag        string
	Streaking                     int
	ContinuousOrderDay            int
	FeeLevel                      int
	Balance                       decimal.Decimal `gorm:"type:decimal(10,2)"`
	FeeSum                        decimal.Decimal `gorm:"type:decimal(10,2)"`
	NoRealNameAlipayWxRechargeSum decimal.Decimal `gorm:"type:decimal(10,2)"`
	RealNameAlipayWxRechargeSum   decimal.Decimal `gorm:"type:decimal(10,2)"`
	IdName                        string
	IdCard                        string
	BankCard                      string
	WithdrawLimitAmount           decimal.Decimal `gorm:"type:decimal(10,2)"`
	RechargeDate                  int
	LianShuoFirstNoticeStatus     string
	LianShuoSecondNoticeStatus    string
	RiskScore                     int
	CanLevelUp                    int
	CustomerLevel                 int
	CurrentExp                    int
	AmountHideFlag                int
	AllowShareOrderInfo           int
	AllowCloseOrderMessage        int
	AllowSystemMessage            int
	AllowAdvertisementPush        int
	IsShowAgricultural            int
	NeedFaceOAuth                 int
	FaceOAuthExpireTime           int
	AlwaysAgreeOrderAgreement     int
	RsaP                          string
	CrpRsa                        string
	RegisterAppToken              string
	RegisterOaid                  string
	RegisterIdfa                  string
	DeviceId                      string
	RegisterAppUuid               string `gorm:"column:registerAppuuid"`
	IsCebBankUpRecordsId          string
	RegisterStartUpRecordsId      string
	SourceSkin                    string
	RiskLevel                     string
	Mute                          int
	Adid                          string `gorm:"column:Adid"`
	Prefsmerchant                 string `gorm:"column:prefsmerchant"`
	FaceOathCertifyId             string
	Theme                         string
}
