package version

import (
	"fmt"
)

var (
	CommitHash     = "n/a"
	BuildTimestamp = "n/a"
)

func BuildVersion() string {
	return fmt.Sprintf("%s (%s)", CommitHash, BuildTimestamp)
}
