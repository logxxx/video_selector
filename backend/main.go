package main

import (
	"github.com/logxxx/utils/runutil"
	"github.com/logxxx/video_selector/base"
	"github.com/logxxx/video_selector/web"
)

func main() {
	base.Init()
	web.InitWeb()

	runutil.WaitForExit()

}
