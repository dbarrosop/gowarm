package core

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"tinygo.org/x/bluetooth"
)

type Central struct {
	adapter     *bluetooth.Adapter
	thermostats map[string]*Thermostat
	logger      *logrus.Entry
}

func New(adapter *bluetooth.Adapter, logger *logrus.Entry, addrs ...string) *Central {
	ths := make(map[string]*Thermostat, len(addrs))
	for _, addr := range addrs {
		ths[addr] = NewThermostat()
	}
	return &Central{
		adapter:     adapter,
		thermostats: ths,
		logger:      logger,
	}
}

func (c *Central) connectionHandler(addr bluetooth.Addresser, connected bool) {
	c.logger.Debugf("connection event for device (%s), connected (%t)", addr, connected)
	for addr, th := range c.thermostats {
		println(addr, th, th.bleDevice)
	}
	if err := c.adapter.StopScan(); err != nil {
		c.logger.Errorf("problem stopping scan: %s", err)
	}

	// 	if connected {
	// 		for a, th := range c.thermostats {
	// 			if a == addr.String() {
	// 				continue
	// 			}
	// 			if th.bleDevice == nil {
	// 				return
	// 			}
	// 		}

	// 		// if we got here we are connected to all devices, so we stop the scan
	// 		if err := c.adapter.StopScan(); err != nil {
	// 			c.logger.Errorf("problem stopping scan: %s", err)
	// 		}
	// 	}

	// 	if !connected {
	// 		c.thermostats[addr.String()].DelDevice()
	// 	}
}

func (c *Central) connectToBLEDevices() error {
	c.logger.Info("scanning for BLE devices")

	err := c.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		for addr, th := range c.thermostats {
			if result.Address.String() == addr {
				c.logger.Infof("device found: %s", result.Address.String())
				device, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
				if err != nil {
					c.logger.Errorf("problem connecting with device: %s", result.Address.String())
				}
				th.SetName(result.LocalName())
				th.SetDevice(device)
				break
			}
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Central) Init() error {
	c.logger.Info("initializing central")

	c.adapter.SetConnectHandler(c.connectionHandler)

	if err := c.connectToBLEDevices(); err != nil {
		return fmt.Errorf("problem connecting to devices: %w", err)
	}

	return nil
}
