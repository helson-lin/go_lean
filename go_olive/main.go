package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// test url: rtsp://rtspstream:aae9ca382fa392561ff0ec3392eb39ec@zephyr.rtsp.stream/movie
// ffmpeg的下载如何集成

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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
	encodedURL := c.Query("url")
	url, err := base64.StdEncoding.DecodeString(encodedURL)
	if err != nil {
		log.Printf("error decoding url: %v", err)
		return
	}
	// 转换的文件类型
	ffmpegOption := []string{}
	var transformType string
	if strings.HasPrefix(string(url), "rtmp") {
		transformType = "rtmp"
		ffmpegOption = append(ffmpegOption, "-rtmp_live", "live")
	} else if strings.HasPrefix(string(url), "rtsp") {
		transformType = "rtsp"
		ffmpegOption = append(ffmpegOption, "-rtsp_transport", "tcp", "-buffer_size", "102400")
	} else {
		transformType = "unknown"
	}
	ffmpegOption = append(ffmpegOption, "-i", string(url), "-acodec", "aac", "-ar", "11025", "-vcodec", "libx264", "-f", "flv")
	// ffmpeg -rtsp_transport tcp -buffer_size 102400 -i rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mp4 -acodec aac -ar 11025 -vcodec libx264 -f flv
	fmt.Println("unknown protocal", transformType, ffmpegOption)
	// 使用ffmpeg获取视频流 左边是赋值 右边是展开
	// 启动异步执行 ffmpeg 命令
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("ffmpeg goroutine panicked: %v", r)
			}
		}()

		cmd := exec.Command("ffmpeg", ffmpegOption...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("cmd stdout pipe error: %v", err)
			return
		}

		if err := cmd.Start(); err != nil {
			log.Printf("cmd start error: %v", err)
			return
		}

		defer func() {
			if err := cmd.Process.Kill(); err != nil {
				log.Printf("Error killing FFmpeg process: %v", err)
			}
		}()

		// 创建一个 goroutine 来处理视频数据并发送到 websocket
		go func() {
			defer stdout.Close()
			buf := make([]byte, 1024)
			for {
				n, err := stdout.Read(buf)
				if err != nil {
					log.Printf("error reading ffmpeg output: %v", err)
					break
				}
				if n > 0 {
					if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
						log.Printf("Error writing to WebSocket: %v", err)
						return
					}
				}
				if err = conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					log.Printf("error writing websocket message: %v", err)
					break
				}
			}
		}()

		if err := cmd.Wait(); err != nil {
			log.Printf("ffmpeg command failed: %v", err)
		}
	}()
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
	r.Use(CORSMiddleware())
	r.Static("/pg", "./static")
	fmt.Println(getFFMpeg())
	// 设置路由
	r.GET("/live/:uid", serveLive)

	// 启动Gin服务器
	port := 8005 // 从配置文件或环境变量读取
	r.Run(fmt.Sprintf(":%d", port))
}
