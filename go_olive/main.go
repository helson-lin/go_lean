package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

// test url: rtsp://rtspstream:aae9ca382fa392561ff0ec3392eb39ec@zephyr.rtsp.stream/movie
// ffmpeg的下载如何集成

type Config struct {
	FFmpeg string `yaml:"ffmpeg"`
	Port   int    `yaml:"port"`
}

var config Config

// readConfig 从配置文件读取配置。
func readConfig() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("config.yml not found, please set config.yml: %v", err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("config.yml is not correct, please check it: %v", err)
	}
	fmt.Printf("FFMPEG ENV Set By Config.yml: %s\n", config.FFmpeg)
}

// validateFFmpeg 验证FFmpeg路径是否有效。
func validateFFmpeg(ffmpegPath string) {
	if _, err := os.Stat(ffmpegPath); os.IsNotExist(err) {
		log.Fatalf("FFmpeg binary specified in config.yml does not exist: %s", ffmpegPath)
	} else if err := exec.Command(ffmpegPath, "-version").Run(); err != nil {
		log.Fatalf("FFmpeg binary specified in config.yml is not executable or invalid: %s", ffmpegPath)
	}
}

func init() {
	readConfig()
	validateFFmpeg(config.FFmpeg)
}

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

func serveLive(c *gin.Context) {
	conn, err := upgradeConnection(c.Writer, c.Request)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	// defer conn.Close()

	url, err := getDecodedURL(c.Query("url"))
	if err != nil {
		log.Printf("error decoding url: %v", err)
		return
	}

	ffmpegOption := determineStreamType(url)
	fmt.Println(ffmpegOption)
	go streamViaFFmpeg(conn, ffmpegOption)
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return upgrader.Upgrade(w, r, nil)
}

func getDecodedURL(encodedURL string) (string, error) {
	url, err := base64.StdEncoding.DecodeString(encodedURL)
	return string(url), err
}

func determineStreamType(url string) []string {
	var ffmpegOption []string
	if strings.HasPrefix(url, "rtmp") {
		ffmpegOption = append(ffmpegOption, "-rtmp_live", "live")
	} else if strings.HasPrefix(url, "rtsp") {
		ffmpegOption = append(ffmpegOption, "-rtsp_transport", "tcp", "-buffer_size", "102400")
	} else {
		return nil
	}
	ffmpegOption = append(ffmpegOption, "-i", url, "-acodec", "aac", "-ar", "11025", "-vcodec", "libx264", "-f", "flv", "pipe:1")
	return ffmpegOption
}

func streamViaFFmpeg(conn *websocket.Conn, ffmpegOption []string) {
	cmd := exec.Command(config.FFmpeg, ffmpegOption...)
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
			log.Printf("error killing FFmpeg process: %v", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("error reading ffmpeg output: %v", err)
			}
			break
		}
		if n > 0 {
			if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				log.Printf("error writing to websocket: %v", err)
				break
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("ffmpeg command failed: %v", err)
	}
}

// var embeddedFiles embed.FS

func main() {
	// gin.SetMode(gin.ReleaseMode)
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.Create("server.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Static("/pg", "./static")
	// staticFiles, _ := fs.Sub(embeddedFiles, "static")
	// r.StaticFS("/pg", http.FS(staticFiles))
	// 设置路由
	r.GET("/live/:uid", serveLive)

	// 启动Gin服务器
	port := config.Port // 从配置文件或环境变量读取
	if port == 0 {
		log.Fatal("Port number is not provided in the config file")
	}

	address := fmt.Sprintf(":%d", port) // 将地址设置为绑定到所有IP地址
	fmt.Printf("Server runing on Port %s\n", address)

	// Start the server
	if err := r.Run(address); err != nil {
		log.Fatalf("Failed to start the server: %v\n", err)
	}
}
