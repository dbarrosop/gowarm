package core

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"tinygo.org/x/bluetooth"
)

type Central struct {
	adapter     *bluetooth.Adapter
	thermostats map[string]*Thermostat
	logger      *logrus.Entry
}

func updateTempValue(addr string) func(float32) {
	return func(v float32) {
		fmt.Printf("%s: %.2fC\n", addr, v)
	}
}

func updateHumidityValue(addr string) func(float32) {
	return func(v float32) {
		fmt.Printf("%s: %.2f%%\n", addr, v)
	}
}

func updateRelayStateValue(addr string) func(bool) {
	return func(v bool) {
		fmt.Printf("%s: is on? %t\n", addr, v)
	}
}

func New(adapter *bluetooth.Adapter, logger *logrus.Entry, addrs ...string) *Central {
	ths := make(map[string]*Thermostat, len(addrs))
	for _, addr := range addrs {
		ths[addr] = NewThermostat(logger, updateTempValue(addr), updateHumidityValue(addr), updateRelayStateValue(addr))
	}
	return &Central{
		adapter:     adapter,
		thermostats: ths,
		logger:      logger,
	}
}

func (c *Central) connectToBLEDevices() error {
	c.logger.Info("scanning for BLE devices")

	err := c.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		for addr, th := range c.thermostats {
			if result.Address.String() == addr {
				// stop early if we are connected already
				if c.thermostats[addr].bleDevice != nil {
					break
				}

				c.logger.Infof("device found: %s", result.Address.String())
				device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
				if err != nil {
					c.logger.Errorf("problem connecting with device: %s", result.Address.String())
				}
				th.SetName(result.LocalName())
				th.SetDevice(device)

				c.logger.Infof("connected to device: %s", th.Name())

				break
			}
		}

		// if we are done connecting to all devices, stop scan
		for _, th := range c.thermostats {
			if th.bleDevice == nil {
				return
			}
		}

		c.logger.Info("stopping scan")
		if err := c.adapter.StopScan(); err != nil {
			c.logger.Warnf("problem stopping scan: %s", err)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Central) Start(ctx context.Context) error {
	c.logger.Info("initializing central")

	if err := c.connectToBLEDevices(); err != nil {
		return fmt.Errorf("problem connecting to devices: %w", err)
	}

	return c.Keepalive()
}

func (c *Central) Keepalive() error {
	d := 15 * time.Second
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			needsReconnect := false

			for _, th := range c.thermostats {
				if time.Since(th.LastSeen) >= 15*time.Second {
					th.logger.Warn("device timed out")
					needsReconnect = true
					if err := th.bleDevice.Disconnect(); err != nil {
						th.logger.Errorf("problem disconnecting device: %s", err)
					}
					th.DelDevice()
				}
			}

			if needsReconnect {
				if err := c.connectToBLEDevices(); err != nil {
					return err
				}
			}

			timer.Reset(d)
		}
	}
}
