package main

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/agent/app"
)

var (
	buildVersion             = "N/A"
	buildDate                = "N/A"
	buildCommit              = "N/A"
	buildInfoMessageTemplate = `Build version: %s
Build date: %s
Build commit: %s
`
)

func main() {
	fmt.Printf(buildInfoMessageTemplate, buildVersion, buildDate, buildCommit)
	app.Run()
}
