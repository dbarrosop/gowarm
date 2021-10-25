package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dbarrosop/gowarm/central/central"
	"github.com/dbarrosop/gowarm/central/homekit"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Core struct {
	central     *central.Central
	homekit     *homekit.Homekit
	thermostats map[string]*Thermostat
	logger      *logrus.Entry
}

func New(c *central.Central, hk *homekit.Homekit, logger *logrus.Entry) *Core {
	return &Core{c, hk, map[string]*Thermostat{}, logger}
}

func (c *Core) AddThermostat(name string, id uint64, address string) {
	th := NewThermostat(c.logger.WithFields(logrus.Fields{"name": name, "address": address, "pkg": "core.thermostat"}))
	c.thermostats[address] = th

	th.ble = c.central.AddThermostat(
		name,
		address,
		th.currentTemperatureCb,
		th.currentHumidityCb,
		th.relayStateCb,
		th.connectCb,
		th.disconnectCb,
	)

	th.hk = c.homekit.AddThermostat(id, address)
}

func (c *Core) InitThermostats() error {
	f := func(address string, ble *central.Thermostat) {
		th := c.thermostats[address]
		th.ble = ble
		c.homekit.InitThermostat(
			address,
			ble.Name(),
			th.targetTemperatureCb,
			th.targetHeatingCoolingStateCb,
		)
	}

	if err := c.central.ConnectToBLEDevices(f); err != nil {
		return err
	}

	return nil
}

func (c *Core) Start(ctx context.Context) error {
	c.logger.Info("starting services")
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		c.logger.Info("starting prometheus handler")

		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", 2114), nil); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error { return c.central.Keepalive(ctx) })

	eg.Go(func() error { return c.homekit.Start() })

	return eg.Wait()
}
