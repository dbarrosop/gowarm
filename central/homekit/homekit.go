package homekit

import (
	"fmt"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/sirupsen/logrus"
)

type Homekit struct {
	thermostats map[string]*Thermostat
	logger      *logrus.Entry
}

func New(logger *logrus.Entry) *Homekit {
	return &Homekit{map[string]*Thermostat{}, logger}
}

func (hk *Homekit) AddThermostat(id uint64, address string) *Thermostat {
	th := NewThermostat(id, hk.logger.WithField("pkg", "homekit.thermostat"))
	hk.thermostats[address] = th
	return th
}

func (hk *Homekit) InitThermostat(address, name string, targetTemperatureCb func(float64), targetHeatingCoolingStateCb func(int)) {
	hk.thermostats[address].Init(name, address, targetTemperatureCb, targetHeatingCoolingStateCb)
}

func (hk *Homekit) Start(name string) error {
	hk.logger.Infof("starting homekit service")

	bridge := accessory.NewBridge(accessory.Info{
		Name:             fmt.Sprintf("GoWarm %s", name),
		SerialNumber:     "66666",
		Manufacturer:     "Potato Industries",
		Model:            "Boiled Potatoes",
		FirmwareRevision: "0.1.0",
		ID:               1,
	})

	accs := make([]*accessory.Accessory, len(hk.thermostats))
	i := 0
	for addr, th := range hk.thermostats {
		hk.logger.WithField("address", addr).Info("adding accessory")
		accs[i] = th.hk.Accessory
		i += 1
	}

	t, err := hc.NewIPTransport(hc.Config{
		StoragePath: "homekit-data",
		Port:        "54114",
		Pin:         "00102003",
		SetupId:     "",
	}, bridge.Accessory, accs...)
	if err != nil {
		return fmt.Errorf("problem creating IP transport: %w", err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	return nil
}
