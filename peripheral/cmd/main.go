// +build tinygo

package main

import (
	"fmt"
	"machine"
	"time"

	"github.com/dbarrosop/gowarm/peripheral/device"
	"github.com/dbarrosop/gowarm/peripheral/thermostat"

	"tinygo.org/x/bluetooth"
)

var (
	name             string
	version          string
	delayLoop                = 2 * time.Second
	hysteresisMargin float32 = 0.1
	targetTemp       float32 = 21.0
	pin                      = machine.D7
)

func bootInfo() {
	time.Sleep(10 * time.Second)
	fmt.Println("Starting gowarm's peripheral")
	fmt.Printf("version: %s\n", version)
	fmt.Printf("name: %s\n", name)
}

func main() {
	bootInfo()

	sensor, err := device.NewBME280Sensor()
	if err != nil {
		println(err.Error())
		return
	}

	relay := device.NewPinRelay(pin)

	th := thermostat.New(sensor, relay, targetTemp, hysteresisMargin)

	ble := device.NewBLE(bluetooth.DefaultAdapter, name, th.SetTargetTemperature, th.SetMode)
	ble.Init()

	prevState := false
	for {
		temp, humidity, state := th.Process()
		fmt.Printf("%.2f C, %.2f %%\n", temp, humidity)

		if err := ble.SendTemperature(temp); err != nil {
			fmt.Printf("problem sending temperature: %s", err)
		}
		if err := ble.SendHumidity(humidity); err != nil {
			fmt.Printf("problem sending humidity: %s", err)
		}

		if prevState != state {
			if err := ble.SendRelayState(state); err != nil {
				fmt.Printf("problem sending relay state: %s", err)
			}
			prevState = state
		}

		time.Sleep(delayLoop)
	}
}
