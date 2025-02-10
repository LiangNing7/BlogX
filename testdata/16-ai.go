package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LiangNing7/BlogX/core"
	"github.com/LiangNing7/BlogX/flags"
	"github.com/LiangNing7/BlogX/global"
)

func chat() {
	url := "https://api.deepseek.com/v1/chat/completions"
	method := "POST"
	payload := strings.NewReader(`{
    "model": "deepseek-chat",
    "messages": [
      {
        "role": "system",
        "content": "你是一个叫凉柠的人工智能助手"
      },
      {
        "role": "user",
        "content": "你好，你是谁?"
      }
    ],
	"stream": false
  }`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

type Choice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason interface{} `json:"finish_reason"`
}
type StreamData struct {
	Id                string      `json:"id"`
	Choices           []Choice    `json:"choices"`
	Created           int         `json:"created"`
	Model             string      `json:"model"`
	Object            string      `json:"object"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

func chatStream() {
	url := "https://api.deepseek.com/v1/chat/completions"
	method := "POST"
	payload := strings.NewReader(`{
    "model": "deepseek-chat",
    "messages": [
      {
        "role": "system",
        "content": "你是一个叫凉柠的人工智能助手"
      },
      {
        "role": "user",
        "content": "你好，你是谁?"
      }
    ],
	"stream": true
  }`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	// 按行分割
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		data := text[6:]
		if data == "[DONE]" {
			break
		}
		var item StreamData
		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			fmt.Printf("解析失败 %s %s", err, data)
			continue
		}
		fmt.Printf(item.Choices[0].Delta.Content)
	}
	fmt.Println()
}

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	chatStream()
}
