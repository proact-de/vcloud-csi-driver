package version

import (
	"runtime"
)

var (
	// String represents the version.
	String = "0.0.0-dev"

	// Revision indicates the commit.
	Revision string

	// Date indicates the build date.
	Date string

	// Go running this binary.
	Go = runtime.Version()
)
