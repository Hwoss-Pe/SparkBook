package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"Webook/bff/web/jwt"
	"Webook/pkg/ginx"
	"Webook/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type GenerateReq struct {
	Content     string `json:"content"`
	Type        string `json:"type"`        // "generate" (default), "polish", "tag"
	Instruction string `json:"instruction"` // For polish: "fix grammar", "expand", etc.
}

type GenerateResp struct {
	Title    string   `json:"title,omitempty"`
	Abstract string   `json:"abstract,omitempty"`
	Content  string   `json:"content,omitempty"` // For polish
	Tags     []string `json:"tags,omitempty"`    // For tag
}

type DifyWorkflowReq struct {
	Inputs       map[string]interface{} `json:"inputs"`
	ResponseMode string                 `json:"response_mode"`
	User         string                 `json:"user"`
}

type DifyWorkflowResp struct {
	Data struct {
		Outputs map[string]interface{} `json:"outputs"`
		Status  string                 `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
	Code    string `json:"code"` // Some error responses use code/message
}

func (a *ArticleHandler) Generate(ctx *gin.Context, req GenerateReq, claims jwt.UserClaims) (ginx.Result, error) {
	if req.Content == "" {
		return ginx.Result{
			Code: 400,
			Msg:  "内容不能为空",
		}, nil
	}

	apiKey := viper.GetString("ai.api_key")
	baseURL := viper.GetString("ai.base_url")

	// 简单的 Mock 逻辑，如果未配置 AI Key，则返回 Mock 数据
	if apiKey == "" || apiKey == "app-xxxx" {
		mockResp := GenerateResp{}
		switch req.Type {
		case "polish":
			mockResp.Content = "这是Mock的润色结果：" + req.Content
		case "tag":
			mockResp.Tags = []string{"Mock标签1", "Mock标签2", "Mock标签3"}
		default:
			mockResp.Title = mockTitle(req.Content)
			mockResp.Abstract = mockAbstract(req.Content)
		}
		return ginx.Result{
			Data: mockResp,
		}, nil
	}

	// 构造 Dify 输入参数
	inputs := map[string]interface{}{
		"content": req.Content,
	}
	if req.Type != "" {
		inputs["type"] = req.Type
	} else {
		// 默认是 generate
		inputs["type"] = "generate"
	}

	if req.Instruction != "" {
		inputs["instruction"] = req.Instruction
	}

	// 增加随机性：在内容末尾添加随机风格指令
	// 仅在生成摘要/标题模式下使用，避免影响润色
	if inputs["type"] == "generate" {
		randomStyles := []string{
			"\n\n(Instruction: Please generate a title that is catchy and intriguing)",
			"\n\n(Instruction: Please generate a title that is professional and concise)",
			"\n\n(Instruction: Please generate a title using a question format)",
			"\n\n(Instruction: Please focus the abstract on the core value proposition)",
			"\n\n(Instruction: Please use a slightly more humorous tone if appropriate)",
			"\n\n(Instruction: Please keep the abstract under 50 words)",
			"", // 也有概率不加任何指令
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		salt := randomStyles[r.Intn(len(randomStyles))]
		inputs["content"] = req.Content + salt
	}

	// 调用 Dify Workflow
	rawOutput, err := a.callDifyWorkflow(ctx, apiKey, baseURL, inputs, claims.Id)
	if err != nil {
		a.l.Error("Dify workflow failed", logger.Error(err))
		return ginx.Result{
			Code: 502,
			Msg:  "AI 生成失败，请稍后重试",
		}, nil
	}

	// 预处理
	cleanOutput := cleanAIOutput(rawOutput)

	// 根据类型解析结果
	var resp GenerateResp
	switch inputs["type"] {
	case "polish":
		// 润色通常直接返回文本
		resp.Content = cleanOutput
	case "tag":
		// 标签解析
		tags, err := parseTags(cleanOutput)
		if err != nil {
			a.l.Warn("Parse tags failed", logger.Error(err))
			// 降级：返回空标签或尝试逗号分割
			resp.Tags = []string{}
		} else {
			resp.Tags = tags
		}
	default:
		// 默认 "generate"
		title, abstract, err := parseGenerateJSON(cleanOutput)
		if err != nil {
			a.l.Warn("Parse generate json failed", logger.Error(err))
			resp.Title = "AI生成标题"
			resp.Abstract = cleanOutput
		} else {
			resp.Title = title
			resp.Abstract = abstract
		}
	}

	return ginx.Result{
		Data: resp,
	}, nil
}

func (a *ArticleHandler) callDifyWorkflow(ctx *gin.Context, apiKey, baseURL string, inputs map[string]interface{}, uid int64) (string, error) {
	url := fmt.Sprintf("%s/workflows/run", strings.TrimRight(baseURL, "/"))

	// 构造请求体
	difyReq := DifyWorkflowReq{
		Inputs:       inputs,
		ResponseMode: "blocking",
		User:         fmt.Sprintf("user-%d", uid),
	}

	jsonData, err := json.Marshal(difyReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("dify api status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var difyResp DifyWorkflowResp
	if err := json.Unmarshal(bodyBytes, &difyResp); err != nil {
		return "", fmt.Errorf("parse dify response failed: %w", err)
	}

	if difyResp.Data.Status != "succeeded" {
		return "", fmt.Errorf("workflow status: %s", difyResp.Data.Status)
	}

	// 假设 Dify Workflow 的 End 节点直接输出 LLM 的结果，变量名通常是 'text' 或自定义的 key
	// 我们遍历 outputs 寻找可能的 JSON 字符串
	var rawOutput string
	if val, ok := difyResp.Data.Outputs["text"]; ok {
		rawOutput = fmt.Sprint(val)
	} else if val, ok := difyResp.Data.Outputs["result"]; ok {
		rawOutput = fmt.Sprint(val)
	} else {
		// 兜底：取第一个 value
		for _, v := range difyResp.Data.Outputs {
			rawOutput = fmt.Sprint(v)
			break
		}
	}

	if rawOutput == "" {
		return "", errors.New("empty output from workflow")
	}

	return rawOutput, nil
}

func cleanAIOutput(raw string) string {
	return strings.TrimSpace(raw)
}

func parseGenerateJSON(raw string) (string, string, error) {
	// 1. 尝试清洗 Markdown 标记
	clean := strings.TrimSpace(raw)
	// 找到第一个 '{'
	start := strings.Index(clean, "{")
	// 找到最后一个 '}'
	end := strings.LastIndex(clean, "}")

	if start == -1 || end == -1 || start >= end {
		// 无法提取 JSON，直接把整个文本作为摘要返回
		return "AI生成标题", clean, nil
	}

	jsonStr := clean[start : end+1]

	var res GenerateResp
	if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
		// 解析失败，降级处理
		return "AI生成标题", clean, nil
	}

	// 简单的默认值处理
	if res.Title == "" {
		res.Title = "AI生成标题"
	}
	if res.Abstract == "" {
		res.Abstract = res.Title
	}

	return res.Title, res.Abstract, nil
}

func parseTags(raw string) ([]string, error) {
	// 尝试解析 JSON 数组 ["tag1", "tag2"]
	clean := strings.TrimSpace(raw)
	start := strings.Index(clean, "[")
	end := strings.LastIndex(clean, "]")
	if start != -1 && end != -1 && start < end {
		jsonStr := clean[start : end+1]
		var tags []string
		if err := json.Unmarshal([]byte(jsonStr), &tags); err == nil {
			return tags, nil
		}
	}

	// 如果不是 JSON 数组，尝试按逗号分割
	tags := strings.Split(clean, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	if len(tags) > 0 {
		return tags, nil
	}

	return nil, errors.New("failed to parse tags")
}

func mockTitle(content string) string {
	runes := []rune(content)
	if len(runes) > 10 {
		return string(runes[:10]) + "..."
	}
	return "示例标题"
}

func mockAbstract(content string) string {
	runes := []rune(content)
	if len(runes) > 50 {
		return string(runes[:50]) + "..."
	}
	return content
}
