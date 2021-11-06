#include <IEEE11073float.h>
#include <unity.h>

void test_decode_temp(uint8_t data[5], float expected) {
    float got = decodeIEEE11073(data + 1);
    TEST_ASSERT_EQUAL_FLOAT(expected, got);
}

void test_encode_temp(float temp, uint8_t expected[5]) {
    uint8_t got[5] = {};
    float2IEEE11073(temp, got + 1);
    TEST_ASSERT_EQUAL_HEX8_ARRAY(expected, got, 5);
}

void test_temp_plus_21point4(void) {
    float temp = 21.4;
    uint8_t bytes[5] = {0x00, 0x60, 0xA7, 0x20, 0xFB};
    test_encode_temp(temp, bytes);
    test_decode_temp(bytes, temp);
}

void test_temp_minus_21point4(void) {
    float temp = -21.4;
    uint8_t bytes[5] = {0x00, 0xA0, 0x58, 0xDF, 0xFB};
    test_encode_temp(temp, bytes);
    test_decode_temp(bytes, temp);
}

void test_temp_plus_21point5(void) {
    float temp = 21.5;
    uint8_t bytes[5] = {0x00, 0xD7, 0x00, 0x00, 0xFF};
    test_encode_temp(temp, bytes);
    test_decode_temp(bytes, temp);
}

void test_temp_minus_21point5(void) {
    float temp = -21.5;
    uint8_t bytes[5] = {0x00, 0x29, 0xFF, 0xFF, 0xFF};
    test_encode_temp(temp, bytes);
    test_decode_temp(bytes, temp);
}

void test_temp_zero(void) {
    float temp = 0.0;
    uint8_t bytes[5] = {0x00, 0x00, 0x00, 0x00, 0x00};
    test_encode_temp(temp, bytes);
    test_decode_temp(bytes, temp);
}

void process() {
    UNITY_BEGIN();
    RUN_TEST(test_temp_plus_21point4);
    RUN_TEST(test_temp_minus_21point4);
    RUN_TEST(test_temp_plus_21point5);
    RUN_TEST(test_temp_minus_21point5);
    RUN_TEST(test_temp_zero);
    UNITY_END();
}

#ifdef ARDUINO

#include <Arduino.h>
void setup() {
    // NOTE!!! Wait for >2 secs
    // if board doesn't support software reset via Serial.DTR/RTS
    delay(2000);

    process();
}

void loop() {
    digitalWrite(13, HIGH);
    delay(100);
    digitalWrite(13, LOW);
    delay(500);
}

#else

int main(int argc, char **argv) {
    process();
    return 0;
}

#endif
