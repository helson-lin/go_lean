package ffmpeg

import (
	"ffandown/utils"
	"fmt"
)

type FFandown struct {
	DIR            string
	FFMPEG_BIN     string
	INPUT_FILE     string
	OUTPUT_FILE    string
	THREAD_NUM     int
	INPUT_OPTIONS  []string
	OUTPUT_OPTIONS []string
}

func (c *FFandown) Log() {
	fmt.Println("ffandown log", c)
}

func (c *FFandown) Init() *FFandown {
	// 设置依赖
	fmt.Println("ffandown initialize")
	localffmpegEnv := utils.GetLocalEnvFFmpeg()
	if localffmpegEnv != "" {
		c.FFMPEG_BIN = localffmpegEnv
	} else {
		libPath, err := utils.InstallFFmpeg("ffmpeg")
		if libPath != "" && err == nil {
			c.FFMPEG_BIN = libPath
		} else {
			fmt.Println("ffmpeg install error", err)
		}
	}
	fmt.Println("ffmpeg", c.FFMPEG_BIN)
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

func (c *FFandown) Run() {
	// 校验输入和输入文件是存在的

}
