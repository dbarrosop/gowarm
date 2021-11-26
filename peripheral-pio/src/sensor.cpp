#include <BlueDot_BME280.h>

namespace Sensor {
BlueDot_BME280 sensor;

void setup() {
    sensor.parameter.communication = 0;         // I2C communication for Sensor 2 (bme2)
    sensor.parameter.I2CAddress = 0x76;         // I2C Address for Sensor 2 (bme2)

    //0b00:     In sleep mode no measurements are performed, but power consumption is at a minimum
    //0b01:     In forced mode a single measured is performed and the device returns automatically to sleep mode
    //0b11:     In normal mode the sensor measures continually (default value)
    sensor.parameter.sensorMode = 0b11;         // Setup Sensor mode for Sensor 2

    sensor.parameter.IIRfilter = 0b000;         // IIR Filter for Sensor 2
    sensor.parameter.humidOversampling = 0b001; // Humidity Oversampling for Sensor 2
    sensor.parameter.tempOversampling = 0b001;  // Temperature Oversampling for Sensor 2
    sensor.parameter.pressOversampling = 0b000; // Pressure Oversampling for Sensor 2

    //0b000:      0.5 msec (default value)
    //0b001:      62.5 msec
    //0b010:      125 msec
    //0b011:      250 msec
    //0b100:      500 msec (0.5 sec, 2Hz)
    //0b101:     1000 msec (1sec, 1Hz)
    //0b110:       10 msec
    //0b111:       20 msec
    sensor.parameter.t_sb = 0b101;

    while (sensor.init() != 0x60 && sensor.init() != 0x58) {
        /* Serial.println(sensor.init()); */
        /* Serial.println("BMx280 sensor not found"); */
        delay(1000);
    }
    /* Serial.println("BMx280 sensor detected"); */
}

float getTemperature() { return sensor.readTempC(); }

float getHumidity() { return sensor.readHumidity(); }

} // namespace Sensor
