package chatgpt

import (
	"sync"
)

type Chat struct {
	Options
	Session struct {
		ConversationId string
		ParentId       string
	}
	mu sync.Mutex
}

type Options struct {
	Headers map[string]string
	BaseURL string
	Model   string
	Retry   int
}

type PartialResponse struct {
	Error error

	ConversationId string      `json:"conversation_id"`
	ResponseError  interface{} `json:"error"`

	Message struct {
		Id     string `json:"id"`
		Status string `json:"status"`

		Author struct {
			Role string `json:"role"`
		} `json:"author"`

		Content struct {
			ContentType string   `json:"content_type"`
			Parts       []string `json:"parts"`
		} `json:"content"`
	} `json:"message"`
}

// 余额
type Billing struct {
	Soft   float64 `json:"soft_limit_usd"`
	Hard   float64 `json:"hard_limit_usd"`
	System float64 `json:"system_hard_limit_usd"`
}
