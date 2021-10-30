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

type Storage interface {
	Close() error
	Write(interface{}) error
}

type Core struct {
	storage     Storage
	central     *central.Central
	homekit     *homekit.Homekit
	Thermostats map[string]*Thermostat
	logger      *logrus.Entry
}

func New(storage Storage, c *central.Central, hk *homekit.Homekit, logger *logrus.Entry) *Core {
	return &Core{
		storage,
		c,
		hk,
		map[string]*Thermostat{},
		logger,
	}
}

func (c *Core) Persist() error {
	return c.storage.Write(c)
}

func (c *Core) AddThermostat(name string, id uint64, address string, config *ThermostatConfig) {
	th := NewThermostat(
		config,
		c.Persist,
		c.logger.WithFields(logrus.Fields{"name": name, "address": address, "pkg": "core.thermostat"}),
	)
	c.Thermostats[address] = th

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
		th := c.Thermostats[address]
		th.ble = ble
		c.homekit.InitThermostat(
			address,
			ble.Name(),
			th.targetTemperatureCb,
			th.targetHeatingCoolingStateCb,
		)

		th.Sync()
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

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", 2114),
			Handler: nil,
		}
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			if err := server.ListenAndServe(); err != nil {
				c.logger.Warning(err)
			}
		}()

		<-ctx.Done()
		return server.Shutdown(context.Background())
	})

	f := func(address string, _ *central.Thermostat) {
		th := c.Thermostats[address]
		th.Sync()
	}

	eg.Go(func() error { return c.central.Keepalive(ctx, f) })

	eg.Go(func() error { return c.homekit.Start() })

	return eg.Wait()
}
