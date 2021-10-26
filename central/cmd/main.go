/*
TODO:
1. Error discovering? gets stuck
2. Get initial state

Peripheral:
1. +-0.1????
2. Recover mechansim
*/
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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

type ThermostatParams struct {
	Thermostats map[string]struct {
		Config *core.ThermostatConfig
	}
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

func stateFile() (core.Storage, *ThermostatParams, error) {
	f, err := os.OpenFile("state.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	tp := &ThermostatParams{
		Thermostats: map[string]struct{ Config *core.ThermostatConfig }{},
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, tp); err != nil {
			return nil, nil, err
		}
	}

	return core.NewFileStorage(f), tp, nil
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

	fs, thparams, err := stateFile()
	if err != nil {
		panic(fmt.Sprintf("problem getting file storage: %s", err))
	}
	defer fs.Close()

	c := core.New(
		fs,
		central.New(adapter, logger.WithField("pkg", "central")),
		homekit.New(logger.WithField("pkg", "hk")),
		logger.WithField("pkg", "core"),
	)

	for _, th := range ths {
		config := &core.ThermostatConfig{
			TargetHeatingCoolingState: 1,
			TargetTemperature:         20.5,
		}
		t, ok := thparams.Thermostats[th.address]
		if ok {
			config = t.Config
		}
		c.AddThermostat(th.room, th.id, th.address, config)
	}

	if err := c.InitThermostats(); err != nil {
		panic(fmt.Sprintf("problem initializing thermostats: %s", err))
	}

	if err := c.Start(context.Background()); err != nil {
		panic(err)
	}
}
