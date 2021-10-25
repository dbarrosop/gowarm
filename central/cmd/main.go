/*
TODO:
1. disconnect gracefully? maybe on the receiving end, test on new
2. Recover state

Peripheral:
1. +-0.2
2. Recover mechansim
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

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

type ThermostatConfig struct {
	room    string
	id      uint64
	address string
}

func parseThermostatConfig() ([]ThermostatConfig, error) {
	thConfig := make([]ThermostatConfig, flag.NArg())
	for i, arg := range flag.Args() {
		a := strings.Split(arg, ",")
		if len(a) != 3 {
			return nil, fmt.Errorf("problem extracting options from arg #%d: %s", i, arg)
		}

		id, err := strconv.Atoi(a[0])
		if err != nil {
			return nil, err
		}

		logrus.Infof("creating configuration entry for %d, %s, %s", id, a[1], a[2])
		thConfig[i] = ThermostatConfig{a[1], uint64(id), a[2]}
	}

	return thConfig, nil
}

func main() {
	flag.Parse()

	ths, err := parseThermostatConfig()
	if err != nil {
		panic(err)
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

	for _, th := range ths {
		c.AddThermostat(th.room, th.id, th.address)
	}

	if err := c.InitThermostats(); err != nil {
		panic(fmt.Sprintf("problem initializing thermostats: %s", err))
	}

	if err := c.Start(context.Background()); err != nil {
		panic(err)
	}
}
