package fund_pay

import "kite/pkg/utils/lian_shuo"

type PostFundPasswordRequest struct {
	parameter lian_shuo.Parameter
}

func (r *PostFundPasswordRequest) Path() string {
	return "/api/fund/oauth/password/create"
}

func (r *PostFundPasswordRequest) WithPassword(password string) *PostFundPasswordRequest {
	r.parameter["password"] = password
	return r
}

func (r *PostFundPasswordRequest) Body() lian_shuo.Parameter {
	return r.parameter
}
