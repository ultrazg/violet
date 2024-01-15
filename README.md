## violet

violet 是一个用于将抖音短视频解析并总结成文本的工具

### 前置

> 本程序依赖 ffmpeg，请确保你的计算机已正确安装

### 配置文件

`config.yaml` 文件，填写必要的配置

```yaml
App:
  # 支持添加多个 Key
  GeminiKey:
    - YOUR-GEMINI-KEY
    - YOUR-GEMINI-KEY
  GeminiUrl: https://generativelanguage.googleapis.com
  UserAgents:
    - Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.2.15
    - Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.66

# 本地服务启动配置项
Server:
  Port: 8080
  Host: localhost

# 代理配置(http|https|socks5://ip:port)，为空则不使用代理
Proxy:
  Protocol: 
```

### 启动

确保配置文件与主程序在同级目录下

```shell
./violet
```

### 接口

#### 去除抖音短视频水印

- 请求方式：`POST`

- 请求地址：`/removeWatermark`

- 请求头：`application/json`

- 请求体参数：

  ```json
  {
      "url": "https://v.douyin.com/iLh7fkpk/"
  }
  ```

- 示例：

  ```shell
  curl --location --request POST 'localhost:8080/removeWatermark' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "url":"https://v.douyin.com/iLh7fkpk/"
  }'
  ```

- 响应：

  ```json
  {
    "message": "success",
    "data": {
      "finalUrl": "去除水印的链接",
      "title": "视频的标题"
    }
  }
  ```

#### 视频转文本

- 请求方式：`POST`

- 请求地址：`/video2text`

- 请求头：`application/json`

- 请求体参数：

  ```json
  {
      "url": "https://v.douyin.com/iLh7fkpk/"
  }
  ```

- 示例：

  ```shell
  curl --location --request POST 'localhost:8080/video2text' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "url":"https://v.douyin.com/iLh7fkpk/"
  }'
  ```

- 响应：

  ```json
  {
    "message": "success",
    "data": {
      "finalUrl": "去除水印的链接",
      "title": "视频的标题",
      "content": "视频的文本"
    }
  }
  ```

### 本地开发

#### 环境

- Go 1.21.3
- ffmpeg

#### 开发

```shell
cd violet
go mod tidy
go run .
```

