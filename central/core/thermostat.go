package core

import (
	"github.com/dbarrosop/gowarm/central/central"
	"github.com/dbarrosop/gowarm/central/homekit"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	currentTempMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "current_temperature",
		Help: "Temperature of the room",
	},
		[]string{"room"},
	)
	currentHumidityMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "current_humidity",
		Help: "Humidity of the room",
	},
		[]string{"room"},
	)
	targetTempMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "target_temperature",
		Help: "Target temperature for the room",
	},
		[]string{"room"},
	)
	relayStateMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "relay_state",
		Help: "Current relay state",
	},
		[]string{"room"},
	)
	connectedMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "connection_state",
		Help: "Current connection state",
	},
		[]string{"room"},
	)
)

type ThermostatConfig struct {
	TargetHeatingCoolingState int
	TargetTemperature         float64
}

type Thermostat struct {
	hk        *homekit.Thermostat
	ble       *central.Thermostat
	persistCb func() error
	Config    *ThermostatConfig
	logger    *logrus.Entry
}

func NewThermostat(config *ThermostatConfig, persistCb func() error, logger *logrus.Entry) *Thermostat {
	logger.Info("creating thermostat")
	return &Thermostat{
		Config:    config,
		persistCb: persistCb,
		logger:    logger,
	}
}

// this method is called when the peripheral is connected
func (th *Thermostat) connectCb() {
	connectedMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(1.0)
}

// this method is called when the peripheral is disconnected
func (th *Thermostat) disconnectCb() {
	connectedMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(0.0)
}

// this method is called by homekit accessory when the target temperature is updated in the home app
func (th *Thermostat) targetTemperatureCb(value float64) {
	th.logger.Infof("changing target temperature to %.2f", value)
	th.ble.SetTargetTemperature(float32(value))
	th.Config.TargetTemperature = value

	if err := th.persistCb(); err != nil {
		th.logger.Errorf("problem trying to persist data: %s", err)
	}

	targetTempMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(value)
}

// this method is called by homekit accessory when the target heating-cooling state is updated in the home app
func (th *Thermostat) targetHeatingCoolingStateCb(value int) {
	th.logger.Infof("changing target heating/cooling state to %d", value)
	th.ble.SetMode([]byte{byte(value)})
	th.Config.TargetHeatingCoolingState = value

	if err := th.persistCb(); err != nil {
		th.logger.Errorf("problem trying to persist data: %s", err)
	}
}

// this method is called by BLE peripheral when the current temperature is updated
func (th *Thermostat) currentTemperatureCb(value float32) {
	th.hk.SetCurrentTemperature(float64(value))

	currentTempMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(float64(value))
}

// this method is called by BLE peripheral when the current humidity is updated
func (th *Thermostat) currentHumidityCb(value float32) {
	currentHumidityMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(float64(value))
}

// this method is called by BLE peripheral when the relay state is updated
func (th *Thermostat) relayStateCb(value bool) {
	s := 0
	if value {
		s = 1
	}

	th.hk.SetCurrentHeatingCoolingState(s)

	relayStateMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(float64(s))
}

func (th *Thermostat) Sync() {
	th.ble.SetTargetTemperature(float32(th.Config.TargetTemperature))
	th.hk.SetTargetTemperature(th.Config.TargetTemperature)

	th.ble.SetMode([]byte{byte(th.Config.TargetHeatingCoolingState)})
	th.hk.SetTargetHeatingCoolingState(th.Config.TargetHeatingCoolingState)

	mode, err := th.ble.GetMode()
	if err != nil {
		th.logger.Errorf("problem getting current mode: %s", err)
	} else {
		th.hk.SetCurrentHeatingCoolingState(int(mode))
	}

	s, err := th.ble.GetRelayState()
	if err != nil {
		th.logger.Errorf("problem getting current relay state: %s", err)
	} else {
		relayStateMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(float64(s))
		th.hk.SetCurrentHeatingCoolingState(int(s))
	}

	targetTempMetric.With(prometheus.Labels{"room": th.ble.Name()}).Set(th.Config.TargetTemperature)
}
