package requests

import "kite/pkg/utils/lian_shuo"

type GetUserInfo struct{}

func NewGetUserInfoRequest() *GetUserInfo {
	return &GetUserInfo{}
}

func (r *GetUserInfo) Path() string {
	return "/api/user/oauth/info/query"
}

func (r *GetUserInfo) Body() lian_shuo.Parameter {
	return nil
}
