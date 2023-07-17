package chatgpt

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

func New(token string, reverseURL string) *Chat {
	if reverseURL == "" {
		reverseURL = BU
	}
	return NewChat(Options{
		Retry:   2,
		BaseURL: reverseURL,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
}

func NewChat(opt Options) *Chat {
	if opt.BaseURL == "" {
		opt.BaseURL = BU
	}

	if opt.Model == "" {
		opt.Model = Gpt3Model
	}

	has := func(key string) bool {
		for k, _ := range opt.Headers {
			if strings.ToLower(k) == key {
				return true
			}
		}
		return false
	}

	for k, v := range H {
		if !has(strings.ToLower(k)) {
			opt.Headers[k] = v
		}
	}

	chat := Chat{Options: opt, Session: struct {
		ConversationId string
		ParentId       string
	}{"", ""}}

	return &chat
}

func (c *Chat) Reply(ctx context.Context, prompt string) (chan PartialResponse, error) {
	var r *http.Response
	c.mu.Lock()

	if c.Retry <= 0 {
		c.Retry = 1
	}

	for index := 1; index <= c.Retry; index++ {
		request, err := c.sendRequest(ctx, prompt)
		if err != nil {
			c.mu.Unlock()
			return nil, NewError(404, err.Error())
		}

		if request.StatusCode == 200 {
			r = request
			break
		}

		if index >= c.Retry {
			c.mu.Unlock()
			_ = request.Body.Close()
			// TODO - more error handle
			return nil, NewError(request.StatusCode, request.Status)
		} else {
			//fmt.Println(err)
		}
	}

	message := make(chan PartialResponse)
	go c.resolve(r, message)
	return message, nil
}

func (c *Chat) sendRequest(ctx context.Context, prompt string) (*http.Response, error) {
	parentId := c.Session.ParentId
	if parentId == "" {
		parentId = uuid.NewString()
	}

	payload := map[string]any{
		"action": "next",
		"model":  c.Model,
		"messages": []map[string]any{
			{
				"id":     uuid.NewString(),
				"author": map[string]string{"role": "user"},
				"content": map[string]any{
					"content_type": "text",
					"parts":        []string{prompt},
				},
			},
		},
		"parent_message_id": parentId,
	}

	if c.Session.ConversationId != "" {
		payload["conversation_id"] = c.Session.ConversationId
		if prompt == "{continue}" {
			payload["action"] = "continue"
			delete(payload, "messages")
		}
	}

	marshal, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/conversation", bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		request.Header.Add(k, v)
	}
	return http.DefaultClient.Do(request)
}

func (c *Chat) resolve(r *http.Response, message chan PartialResponse) {
	defer c.mu.Unlock()
	defer close(message)

	reader := bufio.NewReader(r.Body)
	DONE := []byte("[DONE]")
	var original []byte
	for {
		line, hasMore, err := reader.ReadLine()
		original = append(line)
		if hasMore {
			continue
		}

		if err == io.EOF {
			return
		}

		if err != nil {
			message <- PartialResponse{
				Error: err,
			}
			return
		}

		if len(original) == 0 {
			continue
		}

		if bytes.HasSuffix(original, DONE) {
			return
		}

		dst := make([]byte, len(original))
		copy(dst, original)
		c.originalResolve(dst, message)
		original = make([]byte, 0)
	}
}

func (c *Chat) originalResolve(original []byte, message chan PartialResponse) {
	block := []byte("data: ")
	//fmt.Println("----", string(original), "=")
	if !bytes.HasPrefix(original, block) {
		return
	}

	if !bytes.HasSuffix(original, []byte("}")) {
		return
	}

	original = bytes.TrimPrefix(original, block)
	if bytes.Equal(original, []byte("[DONE]")) {
		return
	}

	var pr PartialResponse
	err := ignorePanicUnmarshal(original, &pr)
	if err != nil {
		//fmt.Println("warn: " + err.Error())
		return
	}

	if pr.Message.Author.Role == "user" {
		return
	}

	if pr.Message.Content.Parts == nil {
		return
	}

	if len(strings.TrimSpace(pr.Message.Content.Parts[0])) == 0 {
		return
	}

	c.Session.ParentId = pr.Message.Id
	if c.Session.ConversationId == "" {
		c.Session.ConversationId = pr.ConversationId
	}
	message <- pr
}

func ignorePanicUnmarshal(data []byte, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("发生了panic:", r)
			if rec, ok := r.(string); ok {
				err = errors.New(rec)
			}
		}
	}()
	return json.Unmarshal(data, v)
}
