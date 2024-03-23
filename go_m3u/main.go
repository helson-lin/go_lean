package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// downloadFile 用于下载文件并保存到本地
func downloadFile(url string, filepath string) error {
	// 需要path内是否存在路径 如果存在需要判断文件夹是否存在 不存在创建文件夹

	// 获取数据
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 写入数据到文件
	_, err = out.ReadFrom(resp.Body)
	return err
}

// parseM3U8 解析m3u8文件并返回.ts文件的链接列表
func parseM3U8(m3u8URL string) ([]string, error) {
	resp, err := http.Get(m3u8URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lines []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, ".ts") {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func main() {
	m3u8URL := "http://devimages.apple.com/iphone/samples/bipbop/gear3/prog_index.m3u8"

	// 解析m3u8文件
	tsFiles, err := parseM3U8(m3u8URL)
	if err != nil {
		panic(err)
	}

	// 下载每个.ts文件
	for i, tsFile := range tsFiles {
		// 如果.ts文件不是完整的URL，需要拼接
		var tsFilePath string
		if !strings.Contains(tsFile, "http://") && !strings.Contains(tsFile, "https://") {
			splitStr := strings.Split(m3u8URL, "/")
			prefixURL := strings.Join(splitStr[:len(splitStr)-1], "/")
			tsFilePath = fmt.Sprintf("%s/%s", prefixURL, tsFile)
			fmt.Println(tsFilePath, tsFile, i)
		}

		filepath := "/temp/" + tsFile // 本地文件名
		fmt.Printf("Downloading %s to %s...\n", tsFilePath, filepath)
		if err := downloadFile(tsFilePath, filepath); err != nil {
			fmt.Printf("Failed to download %s: %v\n", tsFile, err)
		}
	}
}
