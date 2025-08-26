package build

import "fmt"

var (
	buildVersion             = "N/A"
	buildDate                = "N/A"
	buildCommit              = "N/A"
	buildInfoMessageTemplate = `Build version: %s
Build date: %s
Build commit: %s
`
)

// GetBuildInfoMessage возвращает сообщение с информацией о сборке.
func GetBuildInfoMessage() string {
	return fmt.Sprintf(buildInfoMessageTemplate, buildVersion, buildDate, buildCommit)
}
