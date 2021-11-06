#include "ble_ada.h"
#include "sensor.h"

#include <bluefruit.h>

static const char *room = "dev";
static const int delayLoop = 2000;
/* static const int relay_pin = 17; // A3 */
static const int relay_pin = 7; // for itstybity
static const float margin = 0.1;
static const float maxIncrease = 0.5;
static unsigned long recoveryTime = 180000;

float temperatureReading = 0.0;
float humidityReading = 0.0;
float targetTemperature = 20.0;
float resetTemp = 0.0;
bool on = false;
bool active = false;
unsigned long previousIter = 0;
unsigned long lastRecoverAttempt = 0;

void targetTemperatureCb(float temperature) {
    Serial.print("got target temperature: ");
    Serial.println(temperature);
    targetTemperature = temperature;
}

void targetModeCb(int mode) {
    Serial.print("got target mode: ");
    Serial.println(mode);
    active = mode > 0;
}

void setup() {
    while (!Serial)
        delay(10);

    BLE::setup(room, targetTemperatureCb, targetModeCb);
    Sensor::setup();

    pinMode(relay_pin, OUTPUT);
}

void attemptReset() {
    digitalWrite(relay_pin, HIGH);
    delay(500);
    digitalWrite(relay_pin, LOW);
    delay(500);
    digitalWrite(relay_pin, HIGH);
    delay(500);
    digitalWrite(relay_pin, LOW);
}

void loop() {
    unsigned long currentIter = millis();

    temperatureReading = Sensor::getTemperature();
    humidityReading = Sensor::getHumidity();

    Serial.print(temperatureReading);
    Serial.print("C, ");
    Serial.print(humidityReading);
    Serial.println("%");

    if ((!active && on) || (active && on && temperatureReading >= targetTemperature + margin)) {
        Serial.println("turning off relay");
        digitalWrite(relay_pin, LOW);
        on = false;
    } else if (active && !on && temperatureReading <= targetTemperature - margin) {
        Serial.println("turning on relay");
        digitalWrite(relay_pin, HIGH);
        on = true;
    }

    if (temperatureReading > targetTemperature + maxIncrease && active &&
        currentIter - lastRecoverAttempt > recoveryTime && temperatureReading >= resetTemp) {
        Serial.println("reset pin due to overheat!");
        lastRecoverAttempt = currentIter;
        resetTemp = temperatureReading;
        attemptReset();
    } else if (temperatureReading <= targetTemperature + maxIncrease) {
        resetTemp = temperatureReading;
    }

    // if we are disconnected for longer than 15 seconds we reboot
    // we skip this check for the first 5 monutes, in case the central is bookting
    /*         if ( BLEH::disconnectedTime() > 15000 && currentIter - startTime > 300000 ) { */
    /*             Serial.println("reset due to BLE issues!"); */
    /*             PWM::reset(); */
    /*         } */

    BLE::write_temperature(temperatureReading);
    BLE::write_humidity(humidityReading);
    BLE::write_state(on);

    previousIter = currentIter;
    delay(delayLoop);
}
