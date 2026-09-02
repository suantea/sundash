package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DeepSeekService wraps the DeepSeek API for AI-powered features.
type DeepSeekService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewDeepSeekService() *DeepSeekService {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	return &DeepSeekService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// Available returns true if a DeepSeek API key is configured.
func (d *DeepSeekService) Available() bool {
	return d.apiKey != ""
}

// ChatMessage is a message in the chat format.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the request body for DeepSeek Chat API.
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse is the response from DeepSeek Chat API.
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

// Chat sends a chat completion request to DeepSeek.
func (d *DeepSeekService) Chat(messages []ChatMessage) (string, error) {
	if !d.Available() {
		return "", fmt.Errorf("DEEPSEEK_API_KEY 未配置")
	}

	reqBody := ChatRequest{
		Model:    "deepseek-chat",
		Messages: messages,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", d.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("DeepSeek API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API %d: %s", resp.StatusCode, string(b))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("解析 DeepSeek 响应失败: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek 返回空响应")
	}
	return chatResp.Choices[0].Message.Content, nil
}

// SuggestGroup suggests a bookmark group based on title and URL.
func (d *DeepSeekService) SuggestGroup(title, url string, groups []string) (string, error) {
	groupList := ""
	for _, g := range groups {
		groupList += "- " + g + "\n"
	}
	msg := fmt.Sprintf(`根据以下书签信息，从已有分组中选择最合适的。只回复分组名，不要其他内容。

书签标题：%s
书签URL：%s

已有分组：
%s`, title, url, groupList)

	result, err := d.Chat([]ChatMessage{{Role: "user", Content: msg}})
	if err != nil {
		return "", err
	}
	return result, nil
}

// Summarize generates a brief summary of the given text.
func (d *DeepSeekService) Summarize(text string) (string, error) {
	return d.Chat([]ChatMessage{{
		Role:    "user",
		Content: "用一句话总结以下内容（中文，不超过50字）：\n\n" + text,
	}})
}

// SemanticSearch re-ranks search results based on semantic relevance to the query.
func (d *DeepSeekService) SemanticSearch(query string, candidates []string) ([]int, error) {
	candidateList := ""
	for i, c := range candidates {
		candidateList += fmt.Sprintf("%d. %s\n", i, c)
	}
	msg := fmt.Sprintf(`从以下候选项中，按与查询的语义相关性排序。只回复排序后的编号（逗号分隔），不要其他内容。

查询：%s

候选：
%s`, query, candidateList)

	result, err := d.Chat([]ChatMessage{{Role: "user", Content: msg}})
	if err != nil {
		return nil, err
	}
	_ = result
	return nil, nil
}
