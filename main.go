package main

import (
	"github.com/zhiyuan1024/sysinfo/app"
)

func main() {
	appObject := app.NewApplication()
	appObject.Start()
}
