package main

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/agent/app"
	"github.com/mkolibaba/metrics/internal/common/build"
)

func main() {
	fmt.Print(build.GetBuildInfoMessage())
	app.Run()
}
