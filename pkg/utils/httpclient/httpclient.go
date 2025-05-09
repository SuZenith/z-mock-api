package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"kite/pkg/logger"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"
)

var (
	ErrMaxRetriesReached = errors.New("maximum number of retries reached")
	ErrInvalidResponse   = errors.New("invalid response")
)

type BackoffStrategy int

const (
	// BackoffConstant 固定延迟
	BackoffConstant BackoffStrategy = iota
	// BackoffExponential 指数退避
	BackoffExponential
	// BackoffExponentialWithJitter 带抖动的指数退避
	BackoffExponentialWithJitter
)

type RetryPolicy func(resp *http.Response, err error) bool

// DefaultRetryPolicy 默认重试策略
func DefaultRetryPolicy(resp *http.Response, err error) bool {
	if err != nil {
		// 如果是网络错误，则重试
		var netErr net.Error
		if errors.As(err, &netErr) {
			return netErr.Temporary() || netErr.Timeout()
		}
		// 其他错误不重试
		return false
	}
	// 服务器错误以及 429(Too Many Requests) 重试
	return resp.StatusCode >= 500 || resp.StatusCode == 429
}

type RetryConfig struct {
	MaxRetries      int             // 最大重试次数
	InitialDelay    time.Duration   // 初始延迟
	MaxDelay        time.Duration   // 最大延迟
	BackoffStrategy BackoffStrategy // 退避策略
	Multiplier      float64         // 退避因子
	RetryPolicy     RetryPolicy     // 重试策略
}

type ClientConfig struct {
	BaseURL           string        // 基础 URL
	Timeout           time.Duration // 请求超时时间
	RetryConfig       RetryConfig   // 重试配置
	EnableCompression bool          // 是否启用压缩
	DisableKeepAlives bool          // 是否禁用长连接
}

func DefaultConfig() ClientConfig {
	return ClientConfig{
		Timeout: 30 * time.Second,
		RetryConfig: RetryConfig{
			MaxRetries:      3,
			InitialDelay:    100 * time.Millisecond,
			MaxDelay:        30 * time.Second,
			BackoffStrategy: BackoffExponentialWithJitter,
			Multiplier:      2.0,
			RetryPolicy:     DefaultRetryPolicy,
		},
		EnableCompression: true,
		DisableKeepAlives: false,
	}
}

// Client HTTP 客户端
type Client struct {
	client  *http.Client
	config  ClientConfig
	baseUrl string
	rand    *rand.Rand
}

func New(config ClientConfig) *Client {
	defaultCfg := DefaultConfig()
	if config.Timeout == 0 {
		config.Timeout = defaultCfg.Timeout
	}
	// 设置重试配置默认值
	if config.RetryConfig.MaxRetries == 0 {
		config.RetryConfig.MaxRetries = defaultCfg.RetryConfig.MaxRetries
	}
	if config.RetryConfig.InitialDelay == 0 {
		config.RetryConfig.InitialDelay = defaultCfg.RetryConfig.InitialDelay
	}
	if config.RetryConfig.MaxDelay == 0 {
		config.RetryConfig.MaxDelay = defaultCfg.RetryConfig.MaxDelay
	}
	if config.RetryConfig.Multiplier == 0 {
		config.RetryConfig.Multiplier = defaultCfg.RetryConfig.Multiplier
	}
	if config.RetryConfig.RetryPolicy == nil {
		config.RetryConfig.RetryPolicy = defaultCfg.RetryConfig.RetryPolicy
	}
	// 创建 HTTP 传输
	transport := &http.Transport{
		DisableCompression: !config.EnableCompression,
		DisableKeepAlives:  config.DisableKeepAlives,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}
	// 初始化随机数生成器
	source := rand.NewSource(time.Now().UnixNano())

	return &Client{
		client:  client,
		config:  config,
		baseUrl: config.BaseURL,
		rand:    rand.New(source),
	}
}

func (c *Client) calculateBackoff(attempt int) time.Duration {
	retryConfig := c.config.RetryConfig
	initialDelay := retryConfig.InitialDelay
	maxDelay := retryConfig.MaxDelay
	multiplier := retryConfig.Multiplier

	var backoff time.Duration

	switch retryConfig.BackoffStrategy {
	case BackoffConstant:
		// 固定延迟
		backoff = initialDelay
	case BackoffExponential:
		// 指数退避: initialDelay * multiplier^(attempt - 1)
		backoff = time.Duration(float64(initialDelay) * math.Pow(multiplier, float64(attempt-1)))
	case BackoffExponentialWithJitter:
		// 带抖动的指数退避
		base := float64(initialDelay) * math.Pow(multiplier, float64(attempt-1))
		// 添加 0 - 100% 的随机抖动
		jitter := c.rand.Float64() * base
		backoff = time.Duration(base + jitter)
	default:
		// 默认使用指数退避
		backoff = time.Duration(float64(initialDelay) * math.Pow(multiplier, float64(attempt-1)))
	}
	// 确保不超过最大延迟
	if backoff > maxDelay {
		backoff = maxDelay
	}

	return backoff
}

func (c *Client) Request(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	fullUrl := url
	if c.baseUrl != "" {
		fullUrl = fmt.Sprintf("%s%s", c.baseUrl, url)
	}
	// 准备请求体
	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		switch b := body.(type) {
		case io.Reader:
			bodyReader = b
		case []byte:
			bodyBytes = b
			bodyReader = bytes.NewReader(b)
		case string:
			bodyBytes = []byte(b)
			bodyReader = bytes.NewReader(bodyBytes)
		default:
			var err error
			bodyBytes, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body: %v", err)
			}
			bodyReader = bytes.NewReader(bodyBytes)
		}
	}
	// 构建请求
	req, err := http.NewRequest(method, fullUrl, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	// 设置默认 headers
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	// 设置自定义 headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 为重试准备 GetBody 函数
	if bodyBytes != nil {
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}
	// 执行请求
	return c.doWithRetry(req)
}

func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var attempt int

	retryConfig := c.config.RetryConfig
	maxRetries := retryConfig.MaxRetries
	retryPolicy := retryConfig.RetryPolicy

	// 记录请求开始
	logger.Info("HTTP request", zap.String("method", req.Method), zap.String("url", req.URL.String()))

	for attempt = 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := c.calculateBackoff(attempt)

			// 记录重试日志
			logger.Info("Retrying request",
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
				zap.Int("attempt", attempt),
				zap.Duration("backoff", backoff),
			)

			time.Sleep(backoff)

			// 如果请求体是 io.ReadCloser, 需要重制
			if req.Body != nil {
				if req.GetBody != nil {
					var err error
					req.Body, err = req.GetBody()
					if err != nil {
						return nil, fmt.Errorf("failed to read request body: %v", err)
					}
				} else {
					// 如果没有 GetBody 函数，无法重试带有请求体的请求
					return nil, errors.New("cannot retry request with body")
				}
			}
		}
		// 执行请求
		resp, err = c.client.Do(req)

		// 检查是否需要重试
		if err == nil && !retryPolicy(resp, err) {
			// 请求成功不需要重试
			break
		}
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				return nil, err
			}
		}
		// 如果达到最大重试次数，则返回错误
		if attempt == maxRetries {
			if err != nil {
				logger.Error("HTTP request failed after retries",
					zap.String("method", req.Method),
					zap.String("url", req.URL.String()),
					zap.Int("attempt", attempt+1),
					zap.Error(err),
				)
				return nil, fmt.Errorf("%w: %v", ErrMaxRetriesReached, err)
			}
			if resp != nil {
				logger.Error("HTTP request failed after retries",
					zap.String("method", req.Method),
					zap.String("url", req.URL.String()),
					zap.Int("attempt", attempt+1),
					zap.Int("status", resp.StatusCode),
				)
			} else {
				logger.Error("HTTP request failed after retries",
					zap.String("method", req.Method),
					zap.String("url", req.URL.String()),
					zap.Int("attempt", attempt+1),
					zap.String("status", "UNKNOWN"),
				)
			}
			return resp, nil
		}
	}
	// 记录请求结果
	if err != nil {
		logger.Error("HTTP request failed",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Int("attempt", attempt+1),
			zap.Error(err),
		)
	}
	if resp != nil {
		logger.Info("HTTP request succeeded",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Int("status", resp.StatusCode),
			zap.Int("attempt", attempt),
		)
	} else {
		logger.Info("HTTP request succeeded",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Int("attempt", attempt),
			zap.String("status", "UNKNOWN"),
		)
	}
	return resp, nil
}

// Get 发送 GET 请求
func (c *Client) Get(url string, headers map[string]string) (*http.Response, error) {
	return c.Request(http.MethodGet, url, nil, headers)
}

// Post 发送 POST 请求
func (c *Client) Post(url string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(http.MethodPost, url, body, headers)
}

// Put 发送 PUT 请求
func (c *Client) Put(url string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(http.MethodPut, url, body, headers)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(url string, headers map[string]string) (*http.Response, error) {
	return c.Request(http.MethodDelete, url, nil, headers)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(url string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(http.MethodPatch, url, body, headers)
}

// GetJSON 发送 GET 请求并解析 JSON 响应
func (c *Client) GetJSON(url string, headers map[string]string, v interface{}) error {
	resp, err := c.Get(url, headers)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body")
		}
	}(resp.Body)

	return parseJSONResponse(resp, v)
}

// PostJSON 发送 POST 请求并解析 JSON 响应
func (c *Client) PostJSON(url string, body, v interface{}, headers map[string]string) error {
	resp, err := c.Post(url, body, headers)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body")
		}
	}(resp.Body)

	return parseJSONResponse(resp, v)
}

// PutJSON 发送 PUT 请求并解析 JSON 响应
func (c *Client) PutJSON(url string, body, v interface{}, headers map[string]string) error {
	resp, err := c.Put(url, body, headers)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body")
		}
	}(resp.Body)

	return parseJSONResponse(resp, v)
}

// DeleteJSON 发送 DELETE 请求并解析 JSON 响应
func (c *Client) DeleteJSON(url string, v interface{}, headers map[string]string) error {
	resp, err := c.Delete(url, headers)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body")
		}
	}(resp.Body)

	return parseJSONResponse(resp, v)
}

// PatchJSON 发送 PATCH 请求并解析 JSON 响应
func (c *Client) PatchJSON(url string, body, v interface{}, headers map[string]string) error {
	resp, err := c.Patch(url, body, headers)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body")
		}
	}(resp.Body)

	return parseJSONResponse(resp, v)
}

// parseJSONResponse 解析 JSON 响应
func parseJSONResponse(resp *http.Response, v interface{}) error {
	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: status code %d, body: %s", ErrInvalidResponse, resp.StatusCode, string(body))
	}

	// 如果 v 为 nil，不需要解析响应
	if v == nil {
		return nil
	}

	// 解析 JSON 响应
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
