package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ffmpeg的下载如何集成

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有CORS请求
	},
}

func serveLive(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// 解码URL（如果是base64编码的）
	encodedURL := c.Param("url")
	url, err := base64.StdEncoding.DecodeString(encodedURL)
	if err != nil {
		log.Printf("error decoding url: %v", err)
		return
	}

	// 使用ffmpeg获取视频流
	cmd := exec.Command("ffmpeg", "-i", string(url), "-f", "flv", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("cmd stdout pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd start error: %v", err)
	}

	// 创建一个goroutine来处理视频数据并发送到websocket
	go func() {
		defer stdout.Close()
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				log.Printf("error reading ffmpeg output: %v", err)
				return
			}
			if err = conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				log.Printf("error writing websocket message: %v", err)
				return
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatalf("cmd wait error: %v", err)
	}
}

func getFFMpeg() string {
	osType := runtime.GOOS   // 操作系统类型，如windows、darwin、linux
	osArch := runtime.GOARCH // CPU架构，如amd64、386、arm64
	if osType == "darwin" {
		return "osx-64"
	} else if osType == "linux" {
		return osType + osArch
	} else {
		return ""
	}
}

func main() {
	r := gin.Default()
	fmt.Println(getFFMpeg())
	// 设置路由
	r.GET("/live/:url", serveLive)

	// 启动Gin服务器
	port := 8005 // 从配置文件或环境变量读取
	r.Run(fmt.Sprintf(":%d", port))
}
