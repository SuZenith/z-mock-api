package fund_pay

import "kite/pkg/utils/lian_shuo"

type GetBankCard struct {
	parameter lian_shuo.Parameter
}

func NewGetBankCardRequest() *GetBankCard {
	return &GetBankCard{}
}

func (r *GetBankCard) Path() string {
	return "/api/fund-pay/oauth/bank/card"
}

func (r *GetBankCard) WithPayCode(payCode string) *GetBankCard {
	r.parameter["pay_code"] = payCode
	return r
}

func (r *GetBankCard) WithPayType(payType int64) *GetBankCard {
	r.parameter["pay_type"] = payType
	return r
}

func (r *GetBankCard) GetBankCard() lian_shuo.Parameter {
	return nil
}
