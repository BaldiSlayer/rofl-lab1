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

func BuildVersionWithLink() string {
	return fmt.Sprintf("[%s](https://github.com/baldiSlayer/rofl-lab1/commit/%s) (%s)", CommitHash, CommitHash, BuildTimestamp)
}
