typedef void (*ble_target_temperature_callback_t)(float temperature);
typedef void (*ble_target_mode_callback_t)(int mode);

namespace BLE {
void setup(const char *name, ble_target_temperature_callback_t targetTempCb,
           ble_target_mode_callback_t targetModeCb);
void loop();
void write_temperature(float temp);
void write_humidity(float humidity);
void write_state(bool state);
} // namespace BLE
