package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/sirupsen/logrus"
)

type Thermostat struct {
	id     uint64
	logger *logrus.Entry
	hk     *accessory.Thermostat
}

func NewThermostat(id uint64, logger *logrus.Entry) *Thermostat {
	logger.Info("creating thermostat")
	return &Thermostat{id, logger, nil}
}

func (th *Thermostat) Init(
	name string,
	serial_number string,
	targetTemperatureCb func(float64),
	targetHeatingCoolingStateCb func(int),
) {
	th.logger = th.logger.WithField("name", name)
	th.logger.Info("initialiting device")

	acc := accessory.NewThermostat(
		accessory.Info{
			Name:             name,
			SerialNumber:     serial_number,
			Manufacturer:     "Potato Industries",
			Model:            "Boiled Potatoes",
			FirmwareRevision: "0.1.0",
			ID:               th.id,
		},
		20.0, 10.0, 40.0, 0.1,
	)

	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(
		func(value float64) {
			targetTemperatureCb(value)
		},
	)

	acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(
		func(value int) {
			targetHeatingCoolingStateCb(value)
		},
	)

	th.hk = acc
}

func (th *Thermostat) SetCurrentHeatingCoolingState(value int) {
	th.hk.Thermostat.CurrentHeatingCoolingState.SetValue(value)
}

func (th *Thermostat) SetTargetHeatingCoolingState(value int) {
	th.hk.Thermostat.TargetHeatingCoolingState.SetValue(value)
}

func (th *Thermostat) SetCurrentTemperature(value float64) {
	th.hk.Thermostat.CurrentTemperature.SetValue(value)
}

func (th *Thermostat) SetTargetTemperature(value float64) {
	th.hk.Thermostat.TargetTemperature.SetValue(value)
}

func (th *Thermostat) GetCurrentHeatingCoolingState() {
	th.hk.Thermostat.CurrentHeatingCoolingState.GetValue()
}

func (th *Thermostat) GetTargetHeatingCoolingState() {
	th.hk.Thermostat.TargetHeatingCoolingState.GetValue()
}

func (th *Thermostat) GetCurrentTemperature() {
	th.hk.Thermostat.CurrentTemperature.GetValue()
}

func (th *Thermostat) GetTargetTemperature() {
	th.hk.Thermostat.TargetTemperature.GetValue()
}
