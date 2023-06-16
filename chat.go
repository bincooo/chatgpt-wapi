package chatgpt

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func New(token string, reverseURL string) (*Chat, error) {
	if reverseURL == "" {
		reverseURL = BU
	}
	return NewChat(Options{
		BaseURL: reverseURL,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
	})
}

func NewChat(opt Options) (*Chat, error) {
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

	return &chat, nil
}

func (c *Chat) Reply(ctx context.Context, prompt string) (chan PartialResponse, error) {
	c.mu.Lock()
	r, err := c.sendRequest(ctx, prompt)
	if err != nil {
		c.mu.Unlock()
		return nil, NewError(404, err.Error())
	}

	if r.StatusCode != 200 {
		c.mu.Unlock()
		_ = r.Body.Close()
		// TODO -
		return nil, NewError(r.StatusCode, r.Status)
	}

	message := make(chan PartialResponse)
	go c.resolve(ctx, r, message)
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

func (c *Chat) resolve(ctx context.Context, r *http.Response, message chan PartialResponse) {
	defer func() {
		c.mu.Unlock()
		_ = r.Body.Close()
		close(message)
	}()

	for {
		reader := bufio.NewReader(r.Body)
	readline:
		original, _, err := reader.ReadLine()

		if err != nil {
			message <- PartialResponse{
				Error: err,
			}
			return
		}

		block := []byte("data: ")
		if !bytes.HasPrefix(original, block) {
			goto readline
		}

		original = bytes.TrimPrefix(original, block)
		if string(original) == "[DONE]" {
			return
		}

		var pr PartialResponse
		err = json.Unmarshal(original, &pr)
		if err != nil {
			message <- PartialResponse{
				Error: err,
			}
			return
		}

		if pr.Message.Author.Role == "user" {
			continue
		}

		if len(strings.TrimSpace(pr.Message.Content.Parts[0])) == 0 {
			continue
		}

		c.Session.ParentId = pr.Message.Id
		c.Session.ConversationId = pr.ConversationId
		message <- pr
	}
}
