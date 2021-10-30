package device

import (
	"errors"
	"machine"

	"tinygo.org/x/drivers/bme280"
	"tinygo.org/x/drivers/bmp280"
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

type BMP280Sensor struct {
	sensor bmp280.Device
}

func NewBMP280Sensor() (*BMP280Sensor, error) {
	machine.I2C0.Configure(machine.I2CConfig{})

	sensor := bmp280.New(machine.I2C0)
	sensor.Address = 0x76
	sensor.Configure(bmp280.STANDBY_1000MS, bmp280.FILTER_2X, bmp280.SAMPLING_2X, bmp280.SAMPLING_2X, bmp280.MODE_NORMAL)

	connected := sensor.Connected()
	if !connected {
		return nil, errors.New("BMP280 not detected")
	}
	return &BMP280Sensor{sensor}, nil
}

func (s *BMP280Sensor) Temperature() (float32, error) {
	temp, err := s.sensor.ReadTemperature()
	if err != nil {
		return 0, err
	}
	return float32(temp) / 1000, nil
}

func (s *BMP280Sensor) Humidity() (float32, error) {
	return 0.0, nil
}
