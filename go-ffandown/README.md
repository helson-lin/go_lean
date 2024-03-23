## ffmpeg command

- 给视频添加水印
`
ffmpeg -y -i 123.mp4 -vf "drawtext=fontfile=/System/Library/Fonts/PingFang.ttc: text='公众号\:影音探
长':x=10:y=10:fontsize=16:fontcolor=DarkGreen:shadowy=2" 123_txt.mp4
`

