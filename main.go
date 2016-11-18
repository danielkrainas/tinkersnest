package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/cmd/root"
	_ "github.com/danielkrainas/tinkersnest/cmd/serve"
	_ "github.com/danielkrainas/tinkersnest/cmd/version"
	"github.com/danielkrainas/tinkersnest/context"
	_ "github.com/danielkrainas/tinkersnest/storage/driver/inmemory"
)

var appVersion string

const DEFAULT_VERSION = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = DEFAULT_VERSION
	}

	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), appVersion)

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}
