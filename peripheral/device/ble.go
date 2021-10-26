package device

import (
	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/peripheral/pkg/types"
)

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

type (
	updateTargetCb func(float32)
	updateModeCb   func(byte)
)

type BLE struct {
	adapter    *bluetooth.Adapter
	name       string
	chTemp     bluetooth.Characteristic
	chHumidity bluetooth.Characteristic
	chMode     bluetooth.Characteristic
	chState    bluetooth.Characteristic
	chTarget   bluetooth.Characteristic
	targetCb   updateTargetCb
	modeCb     updateModeCb
}

func NewBLE(adapter *bluetooth.Adapter, name string, targetCb updateTargetCb, modeCb updateModeCb) *BLE {
	return &BLE{
		adapter:  adapter,
		name:     name,
		targetCb: targetCb,
		modeCb:   modeCb,
	}
}

func (ble *BLE) Init() error {
	println("starting BLE")
	if err := ble.adapter.Enable(); err != nil {
		return err
	}

	adv := ble.adapter.DefaultAdvertisement()
	if err := adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    ble.name,
		ServiceUUIDs: []bluetooth.UUID{ServiceUUIDEnvironmentalSensing},
	}); err != nil {
		return err
	}

	if err := adv.Start(); err != nil {
		return err
	}

	if err := ble.adapter.AddService(&bluetooth.Service{
		UUID: ServiceUUIDEnvironmentalSensing,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &ble.chTemp,
				UUID:   CharacteristicUUIDTemperatureMeasurement,
				Value:  types.Float32bytes(0),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ble.chHumidity,
				UUID:   CharacteristicUUIDHumidity,
				Value:  types.Float32bytes(0),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ble.chState,
				UUID:   CharacteristicUUIDRelayState,
				Value:  []byte{0, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ble.chMode,
				UUID:   CharacteristicUUIDMode,
				Value:  []byte{0, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(_ bluetooth.Connection, _ int, value []byte) {
					// fmt.Printf("received write event for 'mode' %v\n", value)
					ble.modeCb(value[0])
				},
			},
			{
				Handle: &ble.chTarget,
				UUID:   CharacteristicUUIDTargetTemperature,
				Value:  []byte{0, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(_ bluetooth.Connection, _ int, value []byte) {
					// fmt.Printf("received write event for 'target temperature' %v\n", types.Float32frombytes(value))
					ble.targetCb(types.Float32frombytes(value))
				},
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (ble *BLE) SendTemperature(temp float32) error {
	_, err := ble.chTemp.Write(types.Float32bytes(temp))
	return err
}

func (ble *BLE) SendHumidity(humidity float32) error {
	_, err := ble.chHumidity.Write(types.Float32bytes(humidity))
	return err
}

func (ble *BLE) SendRelayState(state bool) error {
	var b byte
	if state {
		b = 0x1
	}
	_, err := ble.chState.Write([]byte{b})
	return err
}
