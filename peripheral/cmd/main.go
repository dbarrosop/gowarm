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
	name    string
	version string
)

const (
	delayLoop                = 2 * time.Second
	hysteresisMargin float32 = 0.1
	targetTemp       float32 = 21.0
	pin                      = machine.D7
	maxIncrease      float32 = 0.5
	recoveryTime             = 3 * time.Minute
	// recoveryTime = 15 * time.Second
)

func bootInfo() {
	time.Sleep(10 * time.Second)
	fmt.Println("Starting gowarm's peripheral")
	fmt.Printf("version: %s\n", version)
	fmt.Printf("name: %s\n", name)
}

func attemptRecover(relay *device.PinRelay) {
	relay.TurnOn()
	time.Sleep(time.Second)
	relay.TurnOff()
	time.Sleep(time.Second)

	relay.TurnOn()
	time.Sleep(time.Second)
	relay.TurnOff()
}

func getBMESensor() thermostat.Sensor {
	var sensor *device.BME280Sensor
	var err error
	for {
		sensor, err = device.NewBME280Sensor()
		if err == nil {
			break
		}
		println(err.Error())
		time.Sleep(delayLoop)
	}
	return sensor
}

func getBMPSensor() thermostat.Sensor {
	var sensor *device.BMP280Sensor
	var err error
	for {
		sensor, err = device.NewBMP280Sensor()
		if err == nil {
			break
		}
		println(err.Error())
		time.Sleep(delayLoop)
	}
	return sensor
}

func main() {
	bootInfo()

	sensor := getBMESensor()
	// sensor := getBMPSensor()

	relay := device.NewPinRelay(pin)

	th := thermostat.New(sensor, relay, targetTemp, hysteresisMargin)

	ble := device.NewBLE(bluetooth.DefaultAdapter, name, th.SetTargetTemperature, th.SetMode)
	ble.Init()

	var lastRecoverAttempt time.Time
	var resetTemp float32
	prevState := false
	for {
		temp, humidity, state := th.Process()
		fmt.Printf("%.2f, %2.f, %t\n", temp, humidity, state)

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

		if temp > th.TargetTemperature()+maxIncrease && th.ModeOn() && time.Since(lastRecoverAttempt) > recoveryTime && temp >= resetTemp {
			resetTemp = temp

			fmt.Println("resetting")
			lastRecoverAttempt = time.Now()
			attemptRecover(relay)
		} else if temp <= th.TargetTemperature()+maxIncrease {
			resetTemp = temp
		}

		time.Sleep(delayLoop)
	}
}
