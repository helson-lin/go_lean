package main

import (
	ffmpeg "ffandown/lib"
)

func main() {
	ff := &ffmpeg.FFandown{
		DIR: "admin12",
	}
	ff.
		Init().
		SetInputFile("https://123.m3u8").
		SetOutputFile("123.mp4").
		Log()
}
