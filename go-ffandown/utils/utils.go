package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// Includes 检查任意类型切片中是否包含特定的元素
func Includes[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func GetOsArch() string {
	osType := runtime.GOOS   // 获取操作系统类型，如windows, darwin, linux
	osArch := runtime.GOARCH // 获取CPU架构，如amd64,, arm	switch osType
	switch osType {
	case "darwin":
		// 对于darwin系统，目前考虑所有为64位
		return "osx-64"
	case "linux":
		// 对于Linux系统，根据架构进行判断
		if osArch == "arm64" {
			return "linux-arm-64"
		} else if osArch == "arm" {
			return "linux-armel"
		} else if osArch == "amd64" {
			return "linux-64"
		} else {
			// 其他情况默认为32位
			return "linux-32"
		}
	case "windows":
		// 对于Windows系统，根据架构进行判断
		if osArch == "amd64" {
			return "win-64"
		} else {
			// 其他情况默认为32位，主要是386架构
			return "win-32"
		}
	default:
		// 对于其他操作系统类型，简单返回操作系统类型和架构
		return osType + "-" + osArch
	}
}

// 安装ffmpeg依赖
func InstallFFmpeg(component string) (string, error) {
	registry := make(map[string]string, 2)
	registry["github"] = "https://nn.oimi.space/https://github.com/ffbinaries/ffbinaries-prebuilt/releases/download/v4.4.1"
	registry["qiniu"] = "https://pic.kblue.site"
	// 文件后缀
	var libPath string
	baseDir := "dep"
	osType := runtime.GOOS
	if osType == "windows" {
		libPath = baseDir + "/ffmpeg.exe"
	} else {
		libPath = baseDir + "/ffmpeg"
	}
	// 如果libPath存在文件并且有权限直接返回
	isPass, err := IsExecutable(libPath)
	if isPass && err == nil {
		return libPath, nil
	}
	// 不存在直接继续
	fileName := component + "-4.4.1-" + GetOsArch() + ".zip"
	downloadPath := registry["qiniu"] + "/" + fileName
	if err := DownloadFile(fileName, downloadPath); err != nil {
		return "", err
	} else {
		// 下载完毕之后，需要解压文件
		if err := Unzip(fileName, baseDir); err != nil {
			return "", err
		}
		if err := FileAuth(libPath); err != nil {
			// 需要考虑如果是win平台后缀肯定是exe
			return "", err
		}
	}
	return libPath, nil
}

// 获取本机的env内的ffmpeg的地址
func GetLocalEnvFFmpeg() string {
	ffmpegPath := os.Getenv("FFMPEG")
	return ffmpegPath
}

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}
	fmt.Print("\n")
	out.Close()
	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

// 解压ZIP文件到指定目录
func Unzip(zipFile, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录
	os.MkdirAll(destDir, 0755)

	// 遍历ZIP文件中的每个文件/目录
	for _, f := range r.File {
		// 计算文件/目录的目标路径
		fpath := filepath.Join(destDir, f.Name)

		// 如果是目录，则创建目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// 确保文件的目录结构被创建
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// 解压文件到目标路径
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close() // 关闭文件以免泄露
			return err
		}

		// 将文件内容复制到目标文件
		_, err = io.Copy(outFile, rc)

		// 关闭文件
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func FileAuth(fileFullPath string) error {
	fileinfo, err := os.Stat(fileFullPath)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	newPermissions := fileinfo.Mode() | 0111

	if err := os.Chmod(fileFullPath, newPermissions); err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	return nil
}

// 检查文件是否具有可执行权限
func IsExecutable(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	// 获取文件模式位
	mode := fileInfo.Mode()

	// 检查所有用户的可执行权限
	isExec := mode&0111 != 0
	return isExec, nil
}
