package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/context"
	_ "github.com/danielkrainas/tinkersnest/tinkerctl/cmd/create"
	_ "github.com/danielkrainas/tinkersnest/tinkerctl/cmd/describe"
	_ "github.com/danielkrainas/tinkersnest/tinkerctl/cmd/get"
	_ "github.com/danielkrainas/tinkersnest/tinkerctl/cmd/login"
	_ "github.com/danielkrainas/tinkersnest/tinkerctl/cmd/ping"
	"github.com/danielkrainas/tinkersnest/tinkerctl/cmd/root"
	_ "github.com/danielkrainas/tinkersnest/tinkerctl/cmd/version"
)

var appVersion string

const DEFAULT_VERSION = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = DEFAULT_VERSION
	}

	ctx := acontext.WithVersion(acontext.Background(), appVersion)

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
