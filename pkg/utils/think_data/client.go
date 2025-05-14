package think_data

import (
	"context"
	"fmt"
	"kite/pkg/logger"
	"kite/pkg/utils/httpclient"
	"time"
)

const (
	ServerUrl = "https://global-receiver-ta.thinkingdata.cn/sync_json"
	AppId     = "3e4424f3e76847b29c1b9c52d07058ba"
	Debug     = 0
)

const (
	ServerRegisterSuccess             = "server_register_success"
	ServerLoginSuccess                = "server_login_success"
	ServerPlaceOrderSuccess           = "server_place_order_success"
	ServerLimitOrderSuccess           = "server_limit_order_success"
	ServerRechargeSuccess             = "server_recharge_success"
	ServerUserActivationSuccess       = "server_user_activation_success"
	ServerUserRegisterSuccess         = "server_user_register_success"
	ServerUserRechargeFirstDaySuccess = "server_user_recharge_first_day_success"

	EventNameApiLoginSuccess                   = "api_login_success"                       // 登录成功
	EventNameApiLoginFail                      = "api_login_fail"                          // 登录失败
	EventNameApiLoginVerificationCodeSuccess   = "api_login_verification_code_success"     // 登录获取验证码成功
	EventNameApiLoginVerificationCodeFail      = "api_login_verification_code_fail"        // 登录获取验证码失败
	EventNameApiPlaceOrderSuccess              = "api_place_order_success"                 // 下单成功
	EventNameApiPlaceOrderFail                 = "api_place_order_fail"                    // 下单失败
	EventNameApiRechargeSuccess                = "api_recharge_success"                    // 充值成功
	EventNameApiRechargeFail                   = "api_recharge_fail"                       // 充值失败
	EventNameApiWithdrawSuccess                = "api_withdraw_success"                    // 提现成功
	EventNameApiWithdrawFail                   = "api_withdraw_fail"                       // 提现失败
	EventNameApiRegisterSuccess                = "api_register_success"                    // 注册成功
	EventNameApiRegisterRiskControlPass        = "api_register_risk_control_pass"          // 注册风控通过
	EventNameApiRiskControlUserForbidden       = "api_risk_control_user_forbidden"         // 注册风控未通过
	EventNameApiFirstRechargeSuccess           = "api_first_recharge_success"              // 首次充值成功
	EventNameApiNewArrivalFirstDayCouponUsed   = "api_new_arrival_first_day_coupon_used"   // 新人券使用成功
	EventNameApiFirstTimePlaceOrderCashSuccess = "api_first_time_place_order_cash_success" // 首次现金下单成功
	EventNameApiPlaceOrderCashSuccess          = "api_place_order_cash_success"            // 现金下单成功
)

type ThinkDataClient interface {
	Track(ctx context.Context, eventName string, userId *uint, appTk *string, ip *string)
}

type thinkDataClient struct {
	client *httpclient.Client
}

func NewThinkDataClient() ThinkDataClient {
	client := httpclient.New(httpclient.DefaultConfig())
	return &thinkDataClient{client}
}

func (td *thinkDataClient) Track(_ context.Context, eventName string, userId *uint, appTk *string, ip *string) {
	body := map[string]interface{}{}
	body["appid"] = AppId
	body["debug"] = Debug
	data := map[string]interface{}{}
	data["account_id"] = fmt.Sprintf("%d", *userId)
	data["distinct_id"] = fmt.Sprintf("%s", *appTk)
	data["type"] = "track"
	data["ip"] = fmt.Sprintf("%s", *ip)
	data["time"] = time.Now().Format("2006-01-02 15:04:05")
	data["event_name"] = eventName
	properties := map[string]interface{}{}
	properties["channelName"] = ""
	properties["channelCode"] = ""
	properties["versionCode"] = ""
	properties["appVersion"] = ""
	properties["skin"] = ""
	properties["theme"] = ""
	properties["platform"] = ""
	properties["imei"] = ""
	properties["oaid"] = ""
	properties["appuuid"] = ""
	properties["apptk"] = ""
	data["properties"] = properties
	body["data"] = data
	headers := map[string]string{}
	err := td.client.PostJSON(ServerUrl, body, nil, headers)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
