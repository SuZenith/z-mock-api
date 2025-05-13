package lian_shuo

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"kite/pkg/utils/httpclient"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	version       = "1.0.0.0"
	baseUrl       = "http://106.14.98.115/gateway"
	platform      = "server"
	appKey        = "B3A9C8BBE0F94AA2899672879851A2F2"
	appSecret     = "F61BCDC7DBEA4C869660AA5E2973607B"
	memberChannel = "055"
	exchange      = "ZJ"
	// 用来生成随机字符串
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen = len(charset)
)

type LianShuo struct {
	client      *httpclient.Client
	accessToken string
}

type Parameter map[string]interface{}

type Request interface {
	Path() string
	Body() Parameter
}

type Response struct {
	RequestId string      `json:"requestId"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

func NewClient(accessToken string) *LianShuo {
	client := httpclient.New(httpclient.ClientConfig{
		BaseURL:           baseUrl,
		Timeout:           time.Second * 10,
		EnableCompression: true,
		DisableKeepAlives: false,
		RetryConfig: httpclient.RetryConfig{
			MaxRetries:      3,
			InitialDelay:    100 * time.Millisecond,
			MaxDelay:        30 * time.Second,
			BackoffStrategy: httpclient.BackoffExponentialWithJitter,
			Multiplier:      2.0,
			RetryPolicy:     httpclient.DefaultRetryPolicy,
		},
	})
	return &LianShuo{
		client:      client,
		accessToken: accessToken,
	}
}

func (l *LianShuo) Request(req Request) error {
	path := req.Path()
	body := l.buildBody(req)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := l.client.Post(path, body, headers)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", response.Body)

	return nil
}

func (l *LianShuo) buildBody(req Request) Parameter {
	// 初始化参数
	timestampMs := time.Now().UnixMilli()
	parameter := make(Parameter)
	parameter["requestId"] = l.generateRequestId(timestampMs)
	parameter["deviceId"] = l.generateRandomString(16)
	parameter["timestamp"] = timestampMs
	parameter["platform"] = platform
	parameter["appKey"] = appKey
	parameter["version"] = version
	parameter["memberChannel"] = memberChannel
	parameter["exchange"] = exchange
	if strings.TrimSpace(l.accessToken) != "" {
		parameter["accessToken"] = l.accessToken
	}
	parameter["data"] = req.Body()
	// 计算签名
	path := req.Path()
	sign := l.calcSignature(path, parameter, appSecret)
	parameter["sign"] = sign

	return parameter
}

func (l *LianShuo) generateRequestId(timestampMs int64) string {
	return fmt.Sprintf("%s-%d-%s", memberChannel, timestampMs, l.generateRandomString(6))
}

func (l *LianShuo) generateRandomString(length int) string {
	b := make([]byte, length)
	r := make([]byte, length+(length/4))
	// 使用 crypto/rand 填充随机字节
	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}
	for i, rb := range r {
		if i >= length {
			break
		}
		c := int(rb) % charsetLen
		b[i] = charset[c]
	}

	return string(b)
}

func (l *LianShuo) calcSignature(path string, parameter Parameter, appSecret string) string {
	// 用 parameter 中的 data 初始化 arguments
	arguments := make(map[string]interface{})
	if data, ok := parameter["data"].(map[string]interface{}); ok {
		for k, v := range data {
			arguments[k] = v
		}
	}
	// 合并 parameter 到 arguments
	for k, v := range parameter {
		arguments[k] = v
	}
	// 删除 data 字段
	delete(arguments, "data")
	// 获取所有键并排序
	keys := make([]string, 0, len(arguments))
	for k := range arguments {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 构建参数字符串
	var argumentStr string
	for _, k := range keys {
		argumentStr += k + fmt.Sprint(arguments[k])
	}
	// 构建签名字符串并计算签名
	signStr := path + argumentStr + appSecret
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(signStr))
	result := h.Sum(nil)

	return url.QueryEscape(base64.StdEncoding.EncodeToString(result))
}
