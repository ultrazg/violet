package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"violet/config"
	"violet/model"
)

func RandKey(t []string) string {
	rand.Seed(time.Now().UnixNano())

	return t[rand.Intn(len(t))]
}

func NewClient() *http.Client {
	// 开启代理
	if config.Config.Proxy.Protocol != "" {
		proxyURL, err := url.Parse(config.Config.Proxy.Protocol)
		if err != nil {
			log.Println("[Error] - 设置代理出错：", err)
			log.Println("用默认直连")
			client := &http.Client{}

			return client
		}

		log.Println("[Info] - 正在代理请求 Gemini\n代理地址：", proxyURL)

		// 创建一个自定义的 Transport
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		// 使用自定义的 Transport 创建一个 http.Client
		client := &http.Client{
			Transport: transport,
		}

		return client
	} else {
		// 全部直连，不走代理
		client := &http.Client{}

		return client
	}
}

func SetGeminiV(data model.FrameInfo) (error, string) {
	if data.Base64Data == "" {
		return fmt.Errorf("[Error] - Base64 数据为空"), ""
	}

	_url := fmt.Sprintf(config.Config.App.GeminiUrl+"/v1beta/models/gemini-pro-vision:generateContent?key=%s", RandKey(config.Config.App.GeminiKey))
	method := "POST"
	payload := model.GeminiData{
		Contents: []model.Contents{
			{
				Parts: []model.Parts{
					{
						Text: fmt.Sprintf("这个图片是一段视频中第%d段的第%d帧，他的详细内容内容是什么？比如有什么人物，他们在做什么动作，说什么话。这个时候你就是一个视频脚本分析大师，你应该剖析他们原本的剧情或者画面呈现的东西，你应该直接输出告诉我视频这一帧呈现的内容，现在请开始你分析：", data.SegmentIndex, data.FrameIndex),
					},
					{
						InlineData: &model.InlineData{
							MimeType: "image/jpeg",
							Data:     data.Base64Data,
						},
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[Error] - error marshaling payload: %v", err), ""
	}

	client := NewClient()

	request, err := http.NewRequest(method, _url, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("[Error] - 创建请求失败：%v", err), ""
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("[Error] - 发送请求失败：%v", err), ""
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("[Error] - 读取响应失败：%v", err), ""
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("[Error] - status code error: %d %s", response.StatusCode, string(body)), ""
	}

	// 解析响应
	var geminiResponse model.GeminiResponse
	err = json.Unmarshal(body, &geminiResponse)
	if err != nil {
		return fmt.Errorf("[Error] - JSON 解析失败 %v", err), ""
	}

	// 检查 Candidates 切片是否为空
	if len(geminiResponse.Candidates) == 0 {
		return fmt.Errorf("[Error] - Candidates 为空"), ""
	}

	// 检查 Parts 切片是否为空
	if len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("[Error] - Parts 为空"), ""
	}

	frameDescription := fmt.Sprintf("片段 %d 中的第 %d 帧的内容是%s",
		data.SegmentIndex,
		data.FrameIndex,
		geminiResponse.Candidates[0].Content.Parts[0].Text)

	return nil, frameDescription
}

func SetGemini(content string) (error, string) {
	_url := fmt.Sprintf(config.Config.App.GeminiUrl+"/v1beta/models/gemini-pro:generateContent?key=%s", RandKey(config.Config.App.GeminiKey))
	method := "POST"

	payload := model.GeminiPro{
		Contents: []model.GeminiProContent{
			{
				Role: "USER",
				Parts: []model.GeminiProPart{
					{
						Text: fmt.Sprintf("你现在是一个视频脚本整合大师。你的任务是将一系列乱序的视频片段整合成一个完整的故事。每个片段都包含了一系列帧的详细内容描述。由于这些片段是乱序的，你需要先找到第 0 片段的第 0 帧，这是视频的开头。从那里开始，确定每个片段及其帧的正确顺序，然后按照这个顺序来分析整个视频的内容。你的最终目标是输出一个连贯的视频内容脚本，该脚本详细地叙述了视频的全部故事线，包括所有关键的对话、场景和情感转变。请注意，你不需要输出处理的过程，只需要提供视频的完整内容概要。现在开始，请查看以下视频片段及其内容描述，并根据这些信息，从第 0 片段的第 0 帧开始，创建一个完整的视频内容脚本：'%s'", content),
					},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[Error] - error marshaling payload: %v", err), ""
	}

	client := NewClient()

	request, err := http.NewRequest(method, _url, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("[Error] - 创建请求失败：%v", err), ""
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("[Error] - 发送请求失败： %v", err), ""
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("[Error] - error reading response body: %v", err), ""
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("[Error] - status code error: %d %s", response.StatusCode, string(body)), ""
	}

	var geminiResponse model.GeminiResponse
	if err = json.Unmarshal(body, &geminiResponse); err != nil {
		return fmt.Errorf("[Error] - JSON 解析失败 %v", err), ""
	}

	if len(geminiResponse.Candidates) == 0 {
		return fmt.Errorf("[Error] - Candidates 为空"), ""
	}

	if len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("[Error] - Parts 为空"), ""
	}

	return nil, geminiResponse.Candidates[0].Content.Parts[0].Text
}
