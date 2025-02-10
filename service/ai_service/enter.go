package ai_service

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LiangNing7/BlogX/global"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
		PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
	} `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

const baseUrl = "https://api.deepseek.com/v1/chat/completions"

func BaseRequest(r Request) (res *http.Response, err error) {
	method := "POST"
	byteData, _ := json.Marshal(r)
	req, err := http.NewRequest(method, baseUrl, bytes.NewBuffer(byteData))
	if err != nil {
		logrus.Errorf("请求参数失败 %s", err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	return
}

//go:embed chat.prompt
var chatPrompt string

func Chat(content string) (msg string, err error) {
	r := Request{
		Model: "deepseek-chat",
		Messages: []Message{
			{
				Role:    "system",
				Content: chatPrompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	res, err := BaseRequest(r)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(res.Body)
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Errorf("解析失败 %s %s", err, string(body))
		return
	}
	if len(response.Choices) > 0 {
		msg = response.Choices[0].Message.Content
		return
	}
	logrus.Errorf("未获取到数据 %s ", string(body))
	return
}
