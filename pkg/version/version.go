package version

import (
	"fmt"
	"runtime"
)

var (
	// These will be set by build flags
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func GetVersionInfo() string {
	return fmt.Sprintf("how version %s (%s) built on %s with %s",
		Version, GitCommit, BuildDate, runtime.Version())
}
