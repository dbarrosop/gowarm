package device

import (
	"errors"
	"machine"

	"tinygo.org/x/drivers/bme280"
)

type BME280Sensor struct {
	sensor bme280.Device
}

func NewBME280Sensor() (*BME280Sensor, error) {
	machine.I2C0.Configure(machine.I2CConfig{})
	sensor := bme280.New(machine.I2C0)
	sensor.Configure()

	connected := sensor.Connected()
	if !connected {
		return nil, errors.New("BME280 not detected")
	}
	return &BME280Sensor{sensor}, nil
}

func (s *BME280Sensor) Temperature() (float32, error) {
	temp, err := s.sensor.ReadTemperature()
	if err != nil {
		return 0, err
	}
	return float32(temp) / 1000, nil
}

func (s *BME280Sensor) Humidity() (float32, error) {
	humidity, err := s.sensor.ReadHumidity()
	if err != nil {
		return 0, err
	}
	return float32(humidity) / 100, nil
}
