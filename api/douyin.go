package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"violet/config"
	"violet/model"
)

// randomUserAgent 从配置文件中随机返回一个 User-Agent 字符串
func randomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return config.Config.App.UserAgents[rand.Intn(len(config.Config.App.UserAgents))]
}

// ProcessUserInput 正则表达式处理输入的字符
func ProcessUserInput(input string) string {
	linkRegex := regexp.MustCompile(`v\.douyin\.com\/[a-zA-Z0-9]+`)
	idRegex := regexp.MustCompile(`\d{19}`)

	if linkRegex.MatchString(input) {
		return linkRegex.FindString(input)
	} else if idRegex.MatchString(input) {
		return idRegex.FindString(input)
	}

	return ""
}

// ExtractVideoId 提取视频的 id
func ExtractVideoId(link string) string {
	// 先判断链接是否包含协议，如果缺失则补充
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	// 发送请求，获取重定向后的 url
	resp, err := http.Get(link)
	if err != nil {
		log.Println("[Error] - 请求失败：", err)

		return ""
	}

	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	log.Print("[Info] - finalURL： " + finalURL)
	finalURL = resp.Request.URL.String()

	idRegex := regexp.MustCompile(`/video/(\d+)`)
	matches := idRegex.FindStringSubmatch(finalURL)

	if len(matches) > 1 {
		log.Println("[Info] - 已获取到视频 ID：" + matches[1])

		return matches[1]
	}
	return ""
}

func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}

// GetDouYinVideoInfo 获取抖音视频的信息
func GetDouYinVideoInfo(videoId string) (string, string, error) {
	url := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s&a_bogus=64745b2b5bdc4e75b720a9a85b19867a", videoId)
	method := "GET"

	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("[Error] - ", err)

		return "", "", err
	}

	request.Header.Add("User-Agent", randomUserAgent())
	response, err := client.Do(request)
	if err != nil {
		log.Println("[Error] - ", err)

		return "", "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("[Error] - ", err)

		return "", "", err
	}

	if response.StatusCode != 200 {
		log.Println("[Error] - ", string(body))

		return "", "", fmt.Errorf("[Error] - 响应失败")
	}

	var videoInfo model.DYVideoInfo
	err = json.Unmarshal(body, &videoInfo)
	if err != nil {
		log.Println("[Error] - JSON 解析失败", err)

		return "", "", err
	}

	if len(videoInfo.ItemList) > 0 && videoInfo.ItemList[0].Video.PlayAddr.Uri != "" {
		uri := videoInfo.ItemList[0].Video.PlayAddr.Uri
		desc := videoInfo.ItemList[0].Desc
		finalURL := fmt.Sprintf("https://www.iesdouyin.com/aweme/v1/play/?video_id=%s&ratio=1080p&line=0", uri)

		return finalURL, desc, nil
	}

	return "", "", fmt.Errorf("[Error] - 找不到视频")
}

func GetDouYinInfo(link string) (string, string, error) {
	videoIdOrLink := ProcessUserInput(link)
	var videoId string
	if videoIdOrLink != "" {
		if IsNumber(videoIdOrLink) {
			videoId = videoIdOrLink
		} else {
			videoId = ExtractVideoId(videoIdOrLink)
		}
	}

	if len(videoId) == 0 {
		return "", "", fmt.Errorf("[Error] - 找不到视频")
	}

	finalURL, title, err := GetDouYinVideoInfo(videoId)
	if err != nil {
		return "", "", err
	}

	return finalURL, title, nil
}
