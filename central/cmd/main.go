package main

import (
	"context"
	"fmt"

	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/central/core"
	"github.com/sirupsen/logrus"
)

var (
	name    string
	version string
)

func main() {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)

	logger := logrus.NewEntry(l)

	logger.WithFields(logrus.Fields{"name": name, "version": version}).Info("starting gowarm-central")

	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		panic(fmt.Sprintf("problem enabling adapter: %s", err))
	}

	c := core.New(adapter, logger, "F5:AC:32:45:C4:AD")
	if err := c.Start(context.Background()); err != nil {
		panic(fmt.Sprintf("problem initializing central: %s", err))
	}
}
