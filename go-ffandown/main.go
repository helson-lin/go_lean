package main

import (
	ffmpeg "ffandown/lib"
)

func main() {
	ff := &ffmpeg.FFandown{
		DIR: "admin12",
	}
	// options := []string{"-acodec", "aac", "-ar", "11025", "-vcodec", "libx264"}
	ff.
		Init()
	// SetInputFile("http://devimages.apple.com/iphone/samples/bipbop/gear3/prog_index.m3u8").
	// SetOutputFile("123.mp4").
	// SetOutputOptions(options).
	// Run()
}
