package main

import (
	"github.com/noborus/tbln/cmd"
)

// Version represents the version. Overwritten by "git describe --tags --abbrev=0"
var Version = "v0.0.1"

// Revision set "git rev-parse --short HEAD"
var Revision = "HEAD"

func main() {
	cmd.Execute(Version, Revision)
}
