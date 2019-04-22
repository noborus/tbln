package main

import (
	"github.com/noborus/tbln/cmd"
)

// Version represents the version
var Version = "0.0.1"

// Revision set "git rev-parse --short HEAD"
var Revision = "HEAD"

func main() {
	cmd.Execute(Version, Revision)
}
