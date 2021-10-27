package central

import (
	"time"

	"github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/peripheral/pkg/types"
)

type (
	floatCb      func(float32)
	boolCb       func(bool)
	connectionCb func()
)

type Thermostat struct {
	logger               *logrus.Entry
	bleDevice            *bluetooth.Device
	name                 string
	tempCb               floatCb
	humidityCb           floatCb
	relayStateCb         boolCb
	connectedCb          connectionCb
	disconnectedCb       connectionCb
	chCurrentTemperature bluetooth.DeviceCharacteristic
	chCTargetTemperature bluetooth.DeviceCharacteristic
	chCurrentRelayState  bluetooth.DeviceCharacteristic
	chCTargetMode        bluetooth.DeviceCharacteristic
	chCurrentHumidity    bluetooth.DeviceCharacteristic
	LastSeen             time.Time
}

func NewThermostat(name string, tempCb, humidityCb floatCb, relayStateCb boolCb, connectCb, disconnectCb connectionCb, logger *logrus.Entry) *Thermostat {
	logger.Info("creating thermostat")
	return &Thermostat{
		name:           name,
		logger:         logger,
		tempCb:         tempCb,
		humidityCb:     humidityCb,
		relayStateCb:   relayStateCb,
		connectedCb:    connectCb,
		disconnectedCb: disconnectCb,
	}
}

func (th *Thermostat) Name() string {
	return th.name
}

func (th *Thermostat) SetDevice(device *bluetooth.Device) {
	th.bleDevice = device

	th.logger.Info("discovering services/characteristics")
	srvcs, err := device.DiscoverServices(
		[]bluetooth.UUID{
			bluetooth.ServiceUUIDGenericAttribute,
			bluetooth.ServiceUUIDEnvironmentalSensing,
		},
	)
	if err != nil {
		panic(err)
	}

	// buffer to retrieve characteristic data
	buf := make([]byte, 255)

	for _, svc := range srvcs {
		switch svc.UUID() {
		case bluetooth.ServiceUUIDGenericAttribute:
			if err := th.discoverGenericAttribute(svc, buf); err != nil {
				panic(err)
			}
		case bluetooth.ServiceUUIDEnvironmentalSensing:
			if err := th.discoverEnvironmentalSensing(svc, buf); err != nil {
				panic(err)
			}
		}
	}
	th.LastSeen = time.Now()

	th.connectedCb()
}

func (th *Thermostat) discoverGenericAttribute(svc bluetooth.DeviceService, buf []byte) error {
	// TODO
	return nil
}

func (th *Thermostat) discoverEnvironmentalSensing(svc bluetooth.DeviceService, buf []byte) error {
	chs, err := svc.DiscoverCharacteristics(
		[]bluetooth.UUID{
			types.CharacteristicUUIDHumidity,
			types.CharacteristicUUIDTemperatureMeasurement,
			types.CharacteristicUUIDMode,
			types.CharacteristicUUIDRelayState,
			types.CharacteristicUUIDTargetTemperature,
			types.CharacteristicUUIDResetAttempt,
		})
	if err != nil {
		return err
	}

	for _, ch := range chs {
		switch ch.UUID() {
		case types.CharacteristicUUIDTemperatureMeasurement:
			th.chCurrentTemperature = ch

			if err := ch.EnableNotifications(func(b []byte) {
				th.LastSeen = time.Now()
				th.tempCb(types.Float32frombytes(b))
			}); err != nil {
				return err
			}
		case types.CharacteristicUUIDHumidity:
			th.chCurrentHumidity = ch

			if err := ch.EnableNotifications(func(b []byte) {
				th.LastSeen = time.Now()
				th.humidityCb(types.Float32frombytes(b))
			}); err != nil {
				return err
			}
		case types.CharacteristicUUIDRelayState:
			th.chCurrentRelayState = ch

			if err := ch.EnableNotifications(func(b []byte) {
				th.LastSeen = time.Now()
				th.relayStateCb(b[0] > 0x0)
			}); err != nil {
				return err
			}
		case types.CharacteristicUUIDMode:
			th.chCTargetMode = ch
		case types.CharacteristicUUIDTargetTemperature:
			th.chCTargetTemperature = ch
		}
	}

	return nil
}

func (th *Thermostat) DelDevice() {
	th.bleDevice = nil

	th.disconnectedCb()
}

func (th *Thermostat) SetTargetTemperature(value float32) error {
	_, err := th.chCTargetTemperature.WriteWithoutResponse(types.Float32bytes(value))
	return err
}

func (th *Thermostat) SetMode(value []byte) error {
	_, err := th.chCTargetMode.WriteWithoutResponse(value)
	return err
}

func (th *Thermostat) GetMode() (byte, error) {
	b := make([]byte, 1)
	_, err := th.chCTargetMode.Read(b)
	if err != nil {
		return 0.0, nil
	}
	return b[0], nil
}
