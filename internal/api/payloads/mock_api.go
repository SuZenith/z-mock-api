package payloads

import (
	"encoding/json"
)

type Headers struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MockApiPayload struct {
	UserId          string    `json:"user_id" validate:"required"`
	Path            string    `json:"path" validate:"required"`
	Method          string    `json:"method" validate:"required"`
	StatusCode      int16     `json:"status_code" validate:"required"`
	ContentType     string    `json:"content_type" validate:"required"`
	Charset         string    `json:"charset" validate:"required"`
	ResponseHeaders []Headers `json:"headers" validate:"required"`
	ResponseBody    string    `json:"response_body" validate:"required"`
}

// GetHeadersJSON 将 ResponseHeaders 转换为指定格式的 JSON
func (m *MockApiPayload) GetHeadersJSON() (json.RawMessage, error) {
	headersMap := make(map[string]string)

	for _, header := range m.ResponseHeaders {
		headersMap[header.Key] = header.Value
	}
	return json.Marshal(headersMap)
}
