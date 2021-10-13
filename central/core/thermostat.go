package core

import (
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

type Thermostat struct {
	bleDevice *bluetooth.Device
	name      string
}

func NewThermostat() *Thermostat {
	return &Thermostat{}
}

func (th *Thermostat) Name() string {
	return th.name
}

func (th *Thermostat) SetName(name string) {
	th.name = name
}

func (th *Thermostat) SetDevice(device *bluetooth.Device) {
	th.bleDevice = device

	println("discovering services/characteristics")
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDGenericAttribute, bluetooth.ServiceUUIDEnvironmentalSensing})
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
		// 		println("- service", svc.UUID().String())

		// 		chars, err := svc.DiscoverCharacteristics(nil)
		// 		if err != nil {
		// 			println(err)
		// 		}
		// 		for _, char := range chars {
		// 			println("-- characteristic", char.UUID().String())
		// 			n, err := char.Read(buf)
		// 			if err != nil {
		// 				println("    ", err.Error())
		// 			} else {
		// 				println("    data bytes", strconv.Itoa(n))
		// 				println("    value =", string(buf[:n]))
		// 			}
		// 		}
	}
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
			n, err := ch.Read(buf)
			if err != nil {
				return err
			}
			println(n, err, buf[:n], types.Float32frombytes(buf[:n]))
			if err := ch.EnableNotifications(func(b []byte) {
				println(types.Float32frombytes(b))
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (th *Thermostat) DelDevice() {
	th.bleDevice = nil
}
