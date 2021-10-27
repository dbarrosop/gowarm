package device

import (
	"tinygo.org/x/bluetooth"

	"github.com/dbarrosop/gowarm/peripheral/pkg/types"
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
		ServiceUUIDs: []bluetooth.UUID{types.ServiceUUIDEnvironmentalSensing},
	}); err != nil {
		return err
	}

	if err := adv.Start(); err != nil {
		return err
	}

	if err := ble.adapter.AddService(&bluetooth.Service{
		UUID: types.ServiceUUIDEnvironmentalSensing,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &ble.chTemp,
				UUID:   types.CharacteristicUUIDTemperatureMeasurement,
				Value:  types.Float32bytes(0),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ble.chHumidity,
				UUID:   types.CharacteristicUUIDHumidity,
				Value:  types.Float32bytes(0),
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ble.chState,
				UUID:   types.CharacteristicUUIDRelayState,
				Value:  []byte{0, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				Handle: &ble.chMode,
				UUID:   types.CharacteristicUUIDMode,
				Value:  []byte{0, 0},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(_ bluetooth.Connection, _ int, value []byte) {
					// fmt.Printf("received write event for 'mode' %v\n", value)
					ble.modeCb(value[0])
				},
			},
			{
				Handle: &ble.chTarget,
				UUID:   types.CharacteristicUUIDTargetTemperature,
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
