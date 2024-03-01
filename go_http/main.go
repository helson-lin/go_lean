package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 总结： gin的使用有点类似于express的操作： use使用中间件，然后提供了一下返回去获取参数或者进行其他的操作
// 数据存储如何操作： 使用什么第三方的库
// 模块化 如何处理

// StatCost 是一个统计耗时请求耗时的中间件
func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set("name", "小王子")
		// 可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
		// 调用该请求的剩余处理程序
		c.Next()
		// 不调用该请求的剩余处理程序
		// c.Abort()
		// 计算耗时
		cost := time.Since(start)
		log.Println(cost)
	}
}

func indexFunc(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "index", "data": "no data"})
}

func homeFunc(c *gin.Context) {
	var address = c.Query("address")
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"address": address,
	})
}

func getMiddlewareParam(c *gin.Context) {
	name := c.MustGet("name").(string)
	log.Println(name)
	c.JSON(http.StatusOK, gin.H{
		"ok": "1",
	})
}

// 获取路由的路径参数
func getPathFunc(c *gin.Context) {
	username := c.Param("username")
	address := c.Param("address")
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"address":  address,
	})
}

// 获取json的请求体
func jsonFunc(c *gin.Context) {
	// 获取json数据
	b, _ := c.GetRawData()
	// 定义map或者结构体
	var m map[string]interface{}
	// 反序列化
	_ = json.Unmarshal(b, &m)
	c.JSON(http.StatusOK, m)
}

func main() {
	S := gin.Default()
	// 启动静态文件服务
	S.Static("/static", "./static/")
	// 使用请求耗时统计中间件
	S.Use(StatCost())
	// 简单的GET请求
	S.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "服务启动成功"})
	})
	S.GET("/index", indexFunc)
	S.GET("/home", homeFunc)
	S.POST("/json", jsonFunc)
	S.GET("/path/:address/:username", getPathFunc)
	S.GET("/test", getMiddlewareParam)

	err := S.Run(":8080")
	if err != nil {
		fmt.Println("服务器启动失败！")
	}
}
