#include <BlueDot_BME280.h>

namespace Sensor {
BlueDot_BME280 sensor;

void setup() {
    sensor.parameter.communication = 0;         // I2C communication for Sensor 2 (bme2)
    sensor.parameter.I2CAddress = 0x76;         // I2C Address for Sensor 2 (bme2)
    sensor.parameter.sensorMode = 0b11;         // Setup Sensor mode for Sensor 2
    sensor.parameter.IIRfilter = 0b100;         // IIR Filter for Sensor 2
    sensor.parameter.humidOversampling = 0b101; // Humidity Oversampling for Sensor 2
    sensor.parameter.tempOversampling = 0b101;  // Temperature Oversampling for Sensor 2
    sensor.parameter.pressOversampling = 0b000; // Pressure Oversampling for Sensor 2

    while (sensor.init() != 0x60) {
        Serial.println("BME280 sensor not found");
        delay(1000);
    }
    Serial.println("BME280 sensor detected");
}

float getTemperature() { return sensor.readTempC(); }

float getHumidity() { return sensor.readHumidity(); }

} // namespace Sensor
