package main

import (
	"github.com/mirakl/http2back/cli"
)

const (
	AppName = "http2back"
	AppDesc = "Push to backends over HTTP"
)

var (
	AppVersion string
)

func main() {
	if AppVersion == "" {
		AppVersion = "master"
	}

	cli.Process(AppName, AppDesc, AppVersion)
}
