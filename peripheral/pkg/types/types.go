package types

import "tinygo.org/x/bluetooth"

// https://yupana-engineering.com/online-uuid-to-c-array-converter
var (
	ServiceUUIDEnvironmentalSensing = bluetooth.ServiceUUIDEnvironmentalSensing

	CharacteristicUUIDHumidity               = bluetooth.CharacteristicUUIDHumidity
	CharacteristicUUIDTemperatureMeasurement = bluetooth.CharacteristicUUIDTemperatureMeasurement

	// 7512cf1b-3595-4723-b5e4-1e4681660d29
	CharacteristicUUIDTargetTemperature = bluetooth.NewUUID([16]byte{0x75, 0x12, 0xcf, 0x1b, 0x35, 0x95, 0x47, 0x23, 0xb5, 0xe4, 0x1e, 0x46, 0x81, 0x66, 0x0d, 0x29})

	// 5a466ead-b952-4a0f-b750-b988104be49d
	CharacteristicUUIDRelayState = bluetooth.NewUUID([16]byte{0x5a, 0x46, 0x6e, 0xad, 0xb9, 0x52, 0x4a, 0x0f, 0xb7, 0x50, 0xb9, 0x88, 0x10, 0x4b, 0xe4, 0x9d})

	// fbf811de-6b33-4a6f-8efc-fddd0f21086d
	CharacteristicUUIDMode = bluetooth.NewUUID([16]byte{0xfb, 0xf8, 0x11, 0xde, 0x6b, 0x33, 0x4a, 0x6f, 0x8e, 0xfc, 0xfd, 0xdd, 0x0f, 0x21, 0x08, 0x6d})
)
