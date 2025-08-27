package main

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/common/build"
	"github.com/mkolibaba/metrics/internal/server/app"
)

func main() {
	fmt.Print(build.GetBuildInfoMessage())
	app.Run()
}
