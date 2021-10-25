package central

import (
	"time"

	"github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/peripheral/pkg/types"
)

// https://yupana-engineering.com/online-uuid-to-c-array-converter
var (
	// 7512cf1b-3595-4723-b5e4-1e4681660d29
	CharacteristicUUIDTargetTemperature = bluetooth.NewUUID([16]byte{0x75, 0x12, 0xcf, 0x1b, 0x35, 0x95, 0x47, 0x23, 0xb5, 0xe4, 0x1e, 0x46, 0x81, 0x66, 0x0d, 0x29})

	// 5a466ead-b952-4a0f-b750-b988104be49d
	CharacteristicUUIDRelayState = bluetooth.NewUUID([16]byte{0x5a, 0x46, 0x6e, 0xad, 0xb9, 0x52, 0x4a, 0x0f, 0xb7, 0x50, 0xb9, 0x88, 0x10, 0x4b, 0xe4, 0x9d})

	// fbf811de-6b33-4a6f-8efc-fddd0f21086d
	CharacteristicUUIDMode = bluetooth.NewUUID([16]byte{0xfb, 0xf8, 0x11, 0xde, 0x6b, 0x33, 0x4a, 0x6f, 0x8e, 0xfc, 0xfd, 0xdd, 0x0f, 0x21, 0x08, 0x6d})
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
			bluetooth.CharacteristicUUIDHumidity,
			bluetooth.CharacteristicUUIDTemperatureMeasurement,
			CharacteristicUUIDMode,
			CharacteristicUUIDRelayState,
			CharacteristicUUIDTargetTemperature,
		})
	if err != nil {
		return err
	}

	for _, ch := range chs {
		switch ch.UUID() {
		case bluetooth.CharacteristicUUIDTemperatureMeasurement:
			th.chCurrentTemperature = ch

			if err := ch.EnableNotifications(func(b []byte) {
				th.LastSeen = time.Now()
				th.tempCb(types.Float32frombytes(b))
			}); err != nil {
				return err
			}
		case bluetooth.CharacteristicUUIDHumidity:
			th.chCurrentHumidity = ch

			if err := ch.EnableNotifications(func(b []byte) {
				th.LastSeen = time.Now()
				th.humidityCb(types.Float32frombytes(b))
			}); err != nil {
				return err
			}
		case CharacteristicUUIDRelayState:
			th.chCurrentRelayState = ch

			if err := ch.EnableNotifications(func(b []byte) {
				th.LastSeen = time.Now()
				th.relayStateCb(b[0] > 0x0)
			}); err != nil {
				return err
			}
		case CharacteristicUUIDMode:
			th.chCTargetMode = ch
		case CharacteristicUUIDTargetTemperature:
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
