package ffmpeg

import (
	"ffandown/utils"
	"fmt"
	"os"
	"os/exec"
)

type FFandown struct {
	DIR            string
	FFMPEG_BIN     string
	FFPROBE_BIN    string
	INPUT_FILE     string
	OUTPUT_FILE    string
	THREAD_NUM     int
	INPUT_OPTIONS  []string
	OUTPUT_OPTIONS []string
}

func (c *FFandown) Log() {
	fmt.Println("ffandown log", c)
}

func (c *FFandown) InstallDepend(libName string) {
	localffmpegEnv := utils.GetLocalEnvFFmpeg(libName)
	if localffmpegEnv != "" {
		if libName == "FFMPEG" {
			c.FFMPEG_BIN = localffmpegEnv
		} else {
			c.FFPROBE_BIN = localffmpegEnv
		}
	} else {
		libPath, err := utils.InstallFFmpeg(libName)
		if libPath != "" && err == nil {
			if libName == "FFMPEG" {
				c.FFMPEG_BIN = libPath
			} else {
				c.FFPROBE_BIN = libPath
			}
		} else {
			fmt.Println("ffmpeg install error", err)
		}
	}
}

func (c *FFandown) Init() *FFandown {
	// 设置依赖
	fmt.Println("ffandown initialize")
	c.InstallDepend("FFMPEG")
	c.InstallDepend("FFPROBE")
	return c
}

func (c *FFandown) SetInputFile(inputFile string) *FFandown {
	c.INPUT_FILE = inputFile
	return c
}

func (c *FFandown) SetOutputFile(outputFile string) *FFandown {
	c.OUTPUT_FILE = outputFile
	return c
}

func (c *FFandown) SetThreadNum(threadNum int) *FFandown {
	c.THREAD_NUM = threadNum
	return c
}

func (c *FFandown) SetInputOptions(inputOptions []string) *FFandown {
	c.INPUT_OPTIONS = append(c.INPUT_OPTIONS, inputOptions...)
	return c
}

func (c *FFandown) SetOutputOptions(outputOptions []string) *FFandown {
	c.OUTPUT_OPTIONS = append(c.OUTPUT_OPTIONS, outputOptions...)
	return c
}

func (c *FFandown) Run() error {
	if c.FFMPEG_BIN == "" {
		c.Init()
	}
	if c.INPUT_FILE == "" || c.OUTPUT_FILE == "" {
		fmt.Println("input or output file is empty")
	}
	// 构建FFmpeg命令参数
	args := append(c.INPUT_OPTIONS, "-i", c.INPUT_FILE)
	args = append(args, c.OUTPUT_OPTIONS...)
	// 如果存在线程开关需要插入线程
	args = append(args, c.OUTPUT_FILE)
	fmt.Println(args)
	// 创建FFmpeg命令
	cmd := exec.Command(c.FFMPEG_BIN, args...)

	// 获取标准输出和标准错误
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("cmd stdout pipe error: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("cmd stderr pipe error: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd start error: %v", err)
	}

	// 读取输出
	go utils.CopyAndLog(stdout, "FFmpeg Output")
	go utils.CopyAndLog(stderr, "FFmpeg Error")

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %v", err)
	}

	// 校验输出文件是否存在
	if _, err := os.Stat(c.OUTPUT_FILE); os.IsNotExist(err) {
		return fmt.Errorf("output file does not exist")
	}

	return nil
}
