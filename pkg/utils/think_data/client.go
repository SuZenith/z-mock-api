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
