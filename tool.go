package chatgpt

import (
	"context"
	"encoding/json"
	"github.com/acheong08/OpenAIAuth/auth"
	encoder "github.com/samber/go-gpt-3-encoder"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

// 计算prompt的token长度
func CalcTokens(prompt string) int {
	resolver, err := encoder.NewEncoder()
	if err != nil {
		logrus.Error(err)
		return 0
	}
	result, err := resolver.Encode(prompt)
	if err != nil {
		logrus.Error(err)
		return 0
	}
	return len(result)
}

// 尾部截取prompt
func TokensEndSubstr(prompt string, maxToken int) string {
	resolver, err := encoder.NewEncoder()
	if err != nil {
		logrus.Error(err)
		return prompt
	}
	result, err := resolver.Encode(prompt)
	if err != nil {
		logrus.Error(err)
		return prompt
	}

	if l := len(result); l > maxToken {
		result = result[l-maxToken:]
	}
	return resolver.Decode(result)
}

// 头部截取prompt
func TokensStartSubstr(prompt string, maxToken int) string {
	resolver, err := encoder.NewEncoder()
	if err != nil {
		logrus.Error(err)
		return prompt
	}
	result, err := resolver.Encode(prompt)
	if err != nil {
		logrus.Error(err)
		return prompt
	}

	if l := len(result); l > maxToken {
		result = result[:maxToken]
	}
	return resolver.Decode(result)
}

// openai-web获取登陆凭证
func WebLogin(email string, passwd string, proxy string) (string, error) {
	authenticator := auth.NewAuthenticator(email, passwd, proxy)
	if err := authenticator.Begin(); err != nil {
		return "", err.Error
	}
	return authenticator.GetAccessToken(), nil
}

// 查询api余额
func Query(ctx context.Context, token string, proxy string) (*Billing, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, BillingURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	request.Header.Add("Authorization", "Bearer "+token)
	if proxy != "" {
		parser, e := url.Parse(proxy)
		if e != nil {
			return nil, e
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(parser),
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	marshal, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var b Billing
	if e := json.Unmarshal(marshal, &b); e != nil {
		return nil, e
	}
	return &b, err
}
