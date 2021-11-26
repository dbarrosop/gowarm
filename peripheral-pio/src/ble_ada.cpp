// 00FB20A760
// 0060A720FB
#include "ble_ada.h"
#include <IEEE11073float.h>

#include <bluefruit.h>

// 181A
const uint8_t THS_UUID_SERVICE[] = {
    0x1A,
    0x18,
};

// 2A1C
const uint8_t THS_UUID_CHR_TEMPERATURE[] = {
    0x1C,
    0x2A,
};

#define THS_UUID_CHR_HUMIDITY 0x2A6F

// 2A6F
/* const uint8_t THS_UUID_CHR_HUMIDITY[] = { */
/*     0x6F, */
/*     0x2A, */
/* }; */

// 7512cf1b-3595-4723-b5e4-1e4681660d29
const uint8_t THS_UUID_CHR_TARGET_TEMPERATURE[] = {0x29, 0x0d, 0x66, 0x81, 0x46, 0x1e, 0xe4, 0xb5,
                                                   0x23, 0x47, 0x95, 0x35, 0x1b, 0xcf, 0x12, 0x75};

// 5a466ead-b952-4a0f-b750-b988104be49d
const uint8_t THS_UUID_CHR_RELAY_STATE[] = {0x9d, 0xe4, 0x4b, 0x10, 0x88, 0xb9, 0x50, 0xb7,
                                            0x0f, 0x4a, 0x52, 0xb9, 0xad, 0x6e, 0x46, 0x5a};

// fbf811de-6b33-4a6f-8efc-fddd0f21086d
const uint8_t THS_UUID_CHR_TARGET_MODE[] = {0x6d, 0x08, 0x21, 0x0f, 0xdd, 0xfd, 0xfc, 0x8e,
                                            0x6f, 0x4a, 0x33, 0x6b, 0xde, 0x11, 0xf8, 0xfb};

ble_target_temperature_callback_t targetTemperatureCb;
ble_target_mode_callback_t targetModeCb;

BLEService ths(UUID16_SVC_ENVIRONMENTAL_SENSING);
BLECharacteristic thsTemperature(UUID16_CHR_TEMPERATURE_MEASUREMENT);
BLECharacteristic thsHumidity(THS_UUID_CHR_HUMIDITY);

void startAdv(void) { // Advertising packet
    Bluefruit.Advertising.addFlags(BLE_GAP_ADV_FLAGS_LE_ONLY_GENERAL_DISC_MODE);
    Bluefruit.Advertising.addTxPower();

    // Include bleuart 128-bit uuid
    Bluefruit.Advertising.addService(ths);

    // Secondary Scan Response packet (optional)
    // Since there is no room for 'Name' in Advertising packet
    Bluefruit.ScanResponse.addName();

    /* Start Advertising
     * - Enable auto advertising if disconnected
     * - Interval:  fast mode = 20 ms, slow mode = 152.5 ms
     * - Timeout for fast mode is 30 seconds
     * - Start(timeout) with timeout = 0 will advertise forever (until
     * connected)
     *
     * For recommended advertising interval
     * https://developer.apple.com/library/content/qa/qa1931/_index.html
     */
    Bluefruit.Advertising.restartOnDisconnect(true);
    Bluefruit.Advertising.setInterval(3200, 3200); // in unit of 0.625 ms
    Bluefruit.Advertising.setFastTimeout(30);   // number of seconds in fast mode
    Bluefruit.Advertising.start(0);             // 0 = Don't stop advertising after n seconds
}

namespace BLE {
void setup(const char *name, ble_target_temperature_callback_t ttCb,
           ble_target_mode_callback_t tmCb) {
    /* Serial.println("Setting up BLE"); */

    targetTemperatureCb = ttCb;
    targetModeCb = tmCb;

    Bluefruit.begin(1, 0);

    // off Blue LED for lowest power consumption
    Bluefruit.autoConnLed(false);
      
    // Set max power. Accepted values are: -40, -30, -20, -16, -12, -8, -4, 0, 4
    Bluefruit.setTxPower(0);

    /* Bluefruit.Periph.setConnectCallback(connect_callback); */
    /* Bluefruit.Periph.setDisconnectCallback(disconnect_callback); */
    Bluefruit.setName(name);

    ths.begin();

    thsTemperature.setFixedLen(5);
    thsTemperature.setProperties(CHR_PROPS_READ | CHR_PROPS_NOTIFY);
    thsTemperature.setPermission(SECMODE_OPEN, SECMODE_NO_ACCESS);
    thsTemperature.begin();
    /* thsTemperature.writeFloat(temperature); */

    thsHumidity.setFixedLen(5);
    thsHumidity.setProperties(CHR_PROPS_READ | CHR_PROPS_NOTIFY);
    thsHumidity.setPermission(SECMODE_OPEN, SECMODE_NO_ACCESS);
    thsHumidity.begin();
    /* thsHumidity.writeFloat(humidity); */

    /* Serial.println("Setting up the advertising"); */
    startAdv();
}

void write_temperature(float temp) {
    uint8_t data[5] = {bit(0)};
    float2IEEE11073(temp, data + 1);
    thsTemperature.notify(data, sizeof(data));
}

void write_humidity(float humidity) {
    uint8_t data[5] = {bit(0)};
    float2IEEE11073(humidity, data + 1);
    thsHumidity.notify(data, sizeof(data));
}

} // namespace BLE
