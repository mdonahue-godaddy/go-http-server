package main

import (
	"github.com/mdonahue-godaddy/go-http-server/pkg/runner"
)

var (
	// commit, version, and data are injected by the make build
	commit  string
	version string
	date    string
)

func main() {
	runner.Run(version, date, commit)
}
