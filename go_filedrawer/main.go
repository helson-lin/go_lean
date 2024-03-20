package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

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
func main() {
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

	// 启动Gin服务器
	port := 7890 // 从配置文件或环境变量读取
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
