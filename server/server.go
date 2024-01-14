package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"violet/api"
)

func Test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func Video2Text(ctx *gin.Context) {
	var data map[string]string

	err := ctx.ShouldBind(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error parsing JSON",
			"error":   err.Error(),
		})

		return
	}
	if len(data["url"]) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数为空",
			"error":   "error parsing JSON",
		})

		return
	}

	var finalUrl, title string
	if strings.Contains(data["url"], "douyin.com") {
		f, t, err := api.GetDouYinInfo(data["url"])
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "视频链接解析失败，请检查链接是否合法",
				"error":   err.Error(),
			})

			return
		}
		finalUrl = f
		title = t
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "未知的视频链接，请检查链接是否合法",
			"error":   "The URL is wrong, please check",
		})

		return
	}

	err, d := api.VideoSlice(finalUrl)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "未知错误，请重试",
			"error":   err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"url":     finalUrl,
			"title":   title,
			"content": d,
		},
	})

	log.Println("[Info] - done")
}

func RemoveWatermark(ctx *gin.Context) {
	var data map[string]string

	err := ctx.ShouldBind(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error parsing JSON",
			"error":   err.Error(),
		})

		return
	}
	if len(data["url"]) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "参数为空",
			"error":   "error parsing JSON",
		})

		return
	}

	var finalUrl, title string
	if strings.Contains(data["url"], "douyin.com") {
		f, t, err := api.GetDouYinInfo(data["url"])
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "视频链接解析失败，请检查链接是否合法",
				"error":   err.Error(),
			})

			return
		}
		finalUrl = f
		title = t
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "未知的视频链接，请检查链接是否合法",
			"error":   "The URL is wrong, please check",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"url":   finalUrl,
			"title": title,
		},
	})

	log.Println("[Info] - done")
}
