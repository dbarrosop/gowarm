/*
TODO:
1. Keepalive
2. prometheus
3. Recover state
4. cli
*/
package main

import (
	"context"
	"fmt"

	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/central/central"
	"github.com/dbarrosop/gowarm/central/core"
	"github.com/dbarrosop/gowarm/central/homekit"
	"github.com/sirupsen/logrus"
)

var (
	name    string
	version string
)

func main() {
	ths := map[uint64]string{
		20: "F5:AC:32:45:C4:AD",
	}

	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)

	logger := logrus.NewEntry(l)

	logger.WithFields(logrus.Fields{"name": name, "version": version}).Info("starting gowarm-central")

	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		panic(fmt.Sprintf("problem enabling adapter: %s", err))
	}

	c := core.New(
		central.New(adapter, logger.WithField("pkg", "central")),
		homekit.New(logger.WithField("pkg", "hk")),
		logger.WithField("pkg", "core"),
	)

	for id, address := range ths {
		c.AddThermostat(id, address)
	}

	if err := c.InitThermostats(); err != nil {
		panic(fmt.Sprintf("problem initializing thermostats: %s", err))
	}

	if err := c.Start(context.Background()); err != nil {
		panic(err)
	}
}
