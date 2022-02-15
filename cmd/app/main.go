package main

import "github.com/ernur-eskermes/web-video-chat/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
