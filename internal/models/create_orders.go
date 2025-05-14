package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type CreateOrders struct {
	Id                        uint64          `gorm:"column:id;primary_key"`
	OrderNum                  string          `gorm:"column:orderNum;type:varchar(30);not null;comment:'订单号'"`
	ProductID                 string          `gorm:"column:productId;type:varchar(10);comment:'商品ID'"`
	ProductName               string          `gorm:"column:productName;type:varchar(255);comment:'商品名称'"`
	ProductCode               string          `gorm:"column:productCode;type:varchar(255);comment:'商品编号'"`
	ProductPic                string          `gorm:"column:productPic;type:varchar(255);comment:'商品图片'"`
	Contract                  string          `gorm:"column:contract;type:varchar(255);not null;comment:'行情类型'"`
	CreatePrice               decimal.Decimal `gorm:"column:createPrice;type:decimal(10,2);comment:'买入价 行情最新价'"`
	CreateQuoteSeq            uint            `gorm:"column:createQuoteSeq;type:int unsigned;comment:'行情序号'"`
	PlaceOrderMode            string          `gorm:"column:placeOrderMode;type:varchar(255);comment:'下单模式:major(专业) speed(极速) activity(活动)'"`
	SuperProductID            string          `gorm:"column:superProductId;type:varchar(255);comment:'关联的标准SKU'"`
	ExclusiveType             int             `gorm:"column:exclusiveType;type:int(255);comment:'专享类型'"`
	CreateTime                uint64          `gorm:"column:createTime;type:bigint(13) unsigned;comment:'下单时间，时间戳'"`
	ActualCreateTime          uint64          `gorm:"column:actualCreateTime;type:bigint(13) unsigned;comment:'成交时间'"`
	Number                    uint            `gorm:"column:number;type:int unsigned;comment:'数量'"`
	Direction                 uint            `gorm:"column:direction;type:int unsigned;comment:'方向类型，2:涨，1:跌'"`
	Amount                    decimal.Decimal `gorm:"column:amount;type:decimal(10,2);comment:'定金'"`
	Fee                       decimal.Decimal `gorm:"column:fee;type:decimal(10,2);comment:'手续费'"`
	LimitType                 uint            `gorm:"column:limitType;type:int unsigned;comment:'止盈止损类型 1-止盈止损比例 2-止盈止损点数'"`
	TopRatio                  decimal.Decimal `gorm:"column:topRatio;type:decimal(10,2) unsigned;comment:'止盈比例'"`
	BottomRatio               decimal.Decimal `gorm:"column:bottomRatio;type:decimal(10,2) unsigned;comment:'止损比例'"`
	TopPoint                  decimal.Decimal `gorm:"column:topPoint;type:decimal(10,2) unsigned;comment:'止盈点数'"`
	BottomPoint               decimal.Decimal `gorm:"column:bottomPoint;type:decimal(10,2) unsigned;comment:'止损点数'"`
	TopPrice                  decimal.Decimal `gorm:"column:topPrice;type:decimal(10,2);comment:'止盈价'"`
	BottomPrice               decimal.Decimal `gorm:"column:bottomPrice;type:decimal(10,2);comment:'止损价'"`
	OvernightFlag             uint            `gorm:"column:overnightFlag;type:int unsigned;comment:'是否持仓过夜单 0：否，1：是'"`
	OvernightSwitch           int             `gorm:"column:overnightSwitch;type:int(10);comment:'持仓过夜开关 0：关，1：开'"`
	LateFee                   decimal.Decimal `gorm:"column:lateFee;type:decimal(10,2);comment:'每日延期费'"`
	LateFeeFrozen             decimal.Decimal `gorm:"column:lateFeeFrozen;type:decimal(10,2);comment:'预付延期费'"`
	CouponID                  string          `gorm:"column:couponId;type:varchar(255);comment:'优惠券id'"`
	CouponName                string          `gorm:"column:couponName;type:varchar(255);comment:'优惠券名称'"`
	UserID                    uint            `gorm:"column:userId;type:int unsigned;comment:'莲说uid'"`
	CreatedAt                 uint            `gorm:"column:created_at;type:int(13) unsigned"`
	UpdatedAt                 uint            `gorm:"column:updated_at;type:int(13) unsigned"`
	OrderType                 int             `gorm:"column:orderType;type:int(10);comment:'订单类型 3-手动退订 4-止盈退订 5-止损退订 6-爆仓退订 7-休市退订'"`
	ClosePrice                decimal.Decimal `gorm:"column:closePrice;type:decimal(10,2);comment:'退订价格'"`
	CloseTime                 int64           `gorm:"column:closeTime;type:bigint(13);comment:'退订时间'"`
	PlAmount                  decimal.Decimal `gorm:"column:plAmount;type:decimal(10,2);default:0.00;comment:'盈亏金额'"`
	ActualClosePrice          decimal.Decimal `gorm:"column:actualClosePrice;type:decimal(10,2);comment:'实际退订行情价格'"`
	CouponFlag                uint            `gorm:"column:couponFlag;type:int(1) unsigned;comment:'是否用券'"`
	FollowedOrdersCount       uint            `gorm:"column:followedOrdersCount;type:int unsigned;comment:'被跟单数量'"`
	FollowUpFlag              uint            `gorm:"column:followUpFlag;type:int(1) unsigned;comment:'是否跟单 1-跟单 0-非跟单'"`
	LikesCount                uint            `gorm:"column:likesCount;type:int unsigned;comment:'点赞数'"`
	ReturnToFollowerPointsSum uint            `gorm:"column:returnToFollowerPointsSum;type:int unsigned;comment:'返回积分数'"`
	CloseAt                   *time.Time      `gorm:"column:closeAt;type:timestamp;comment:'退订时间-bullet'"`
	DeliveryNumber            uint            `gorm:"column:deliveryNumber;type:int unsigned;default:0;comment:'提货数量'"`
	MinLossLimitRate          decimal.Decimal `gorm:"column:minLossLimitRate;type:decimal(10,2) unsigned;comment:'最小止跌比例'"`
	MaxLossLimitRate          decimal.Decimal `gorm:"column:maxLossLimitRate;type:decimal(10,2) unsigned;comment:'最大止跌比例'"`
	MinProfitLimitRate        decimal.Decimal `gorm:"column:minProfitLimitRate;type:decimal(10,2) unsigned;comment:'最小止盈比例'"`
	MaxProfitLimitRate        decimal.Decimal `gorm:"column:maxProfitLimitRate;type:decimal(10,2) unsigned;comment:'最大止盈比例'"`
	CloseTimestamp            int64           `gorm:"column:closeTimestamp;type:bigint(13);comment:'退订时间-bullet'"`
	SupportNumber             string          `gorm:"column:supportNumber;type:varchar(20);comment:'保单号'"`
	OrderNumber               string          `gorm:"column:orderNumber;type:varchar(20);comment:'订单号'"`
	IsDoOrder                 int             `gorm:"column:isDoOrder;type:int(1);comment:'挂单是否下单'"`
	IsLimitOrder              int             `gorm:"column:isLimitOrder;type:int(1);comment:'订单是否为挂单'"`
	LimitPrice                decimal.Decimal `gorm:"column:limitPrice;type:decimal(10,2);comment:'挂单价'"`
	FloatPoint                decimal.Decimal `gorm:"column:floatPoint;type:decimal(10,2);comment:'浮动点数'"`
	OrderStatus               int             `gorm:"column:orderStatus;type:int(1)"`
	Reason                    string          `gorm:"column:reason;type:varchar(255)"`
	SourceSkin                string          `gorm:"column:sourceSkin;type:varchar(255);comment:'记录用户下单的app skin'"`
	RelatedRequestID          string          `gorm:"column:relatedRequestId;type:varchar(36)"`
	Theme                     string          `gorm:"column:theme;type:varchar(255)"`
}
