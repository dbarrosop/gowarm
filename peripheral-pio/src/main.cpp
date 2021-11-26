#include "ble_ada.h"
#include "sensor.h"

#include <bluefruit.h>

static const char *room = "guest room";
static const int delayLoop = 5000;
/* static const int relay_pin = 17; // A3 */
static const int relay_pin = 7; // for itstybity
static const float margin = 0.1;

float temperatureReading = 0.0;
float humidityReading = 0.0;
float targetTemperature = 20.0;
float resetTemp = 0.0;
bool on = false;
bool active = false;
unsigned long previousIter = 0;
unsigned long lastRecoverAttempt = 0;


void targetTemperatureCb(float temperature) {
    /* Serial.print("got target temperature: "); */
    /* Serial.println(temperature); */
    targetTemperature = temperature;
}

void targetModeCb(int mode) {
    /* Serial.print("got target mode: "); */
    /* Serial.println(mode); */
    active = mode > 0;
}

void loop() {
    temperatureReading = Sensor::getTemperature();
    humidityReading = Sensor::getHumidity();

    /* Serial.print(temperatureReading); */
    /* Serial.print("C, "); */
    /* Serial.print(humidityReading); */
    /* Serial.println("%"); */

    BLE::write_temperature(temperatureReading);
    BLE::write_humidity(humidityReading);

    delay(delayLoop);
}

void setup() {
    /* while (!Serial) */
    /*     delay(10); */
    //Enable low power mode
    sd_power_mode_set(NRF_POWER_MODE_LOWPWR);


    Sensor::setup();
    BLE::setup(room, targetTemperatureCb, targetModeCb);
}
