package types

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	MDER_FLOAT_MANTISSA_MAX = 8388605 // 0x007FFFFD
	MDER_FLOAT_EXPONENT_MAX = 127
	MDER_FLOAT_EXPONENT_MIN = -128
	MDER_FLOAT_PRECISION    = 10000000
)

var (
	MDER_POSITIVE_INFINITY = []byte{0x00, 0x00, 0x7F, 0xFF, 0xFE}
	MDER_NEGATIVE_INFINITY = []byte{0x00, 0x00, 0x80, 0x00, 0x02}
)

func DecodeIEEE11073(bytes []byte) float32 {
	b0 := uint8(bytes[1])
	b1 := uint8(bytes[2])
	b2 := uint8(bytes[3])

	exponent := int8(bytes[4])

	// means bit sign is set, so it is a negative number
	sign := byte(0)
	if b2 >= 127 {
		sign = byte(255)
	}

	mantissa := int32(binary.BigEndian.Uint32([]byte{sign, b2, b1, b0}))
	fmt.Println(mantissa)
	f := float64(mantissa) * math.Pow(10.0, float64(exponent))
	return float32(math.Round(float64(f)*100) / 100)
}

func EncodeIEEE11073(value float32) []byte {
	result := make([]byte, 5)
	if value == 0 {
		return result
	}

	sign := 1
	if value < 0 {
		sign = -1
	}

	mantissa := math.Abs(float64(value))
	exponent := 0 // Note: 10**x exponent, not 2**x

	// scale up if number is too big
	for mantissa > MDER_FLOAT_MANTISSA_MAX {
		mantissa /= 10.0
		exponent += 1
		if exponent > MDER_FLOAT_EXPONENT_MAX {
			// argh, should not happen
			if sign < 0 {
				return MDER_NEGATIVE_INFINITY
			}
			return MDER_POSITIVE_INFINITY
		}
	}

	// scale down if number is too small
	for mantissa < 1 {
		mantissa *= 10
		exponent -= 1
		if exponent < MDER_FLOAT_EXPONENT_MIN {
			// argh, should not happen
			panic("no idea what happened")
		}
	}

	// scale down if number needs more precision
	smantissa := math.Round(mantissa * MDER_FLOAT_PRECISION)
	rmantissa := math.Round(mantissa) * MDER_FLOAT_PRECISION
	mdiff := math.Abs(smantissa - rmantissa)
	for mdiff > 0.5 && exponent > MDER_FLOAT_EXPONENT_MIN && (mantissa*10) <= MDER_FLOAT_MANTISSA_MAX {
		mantissa *= 10
		exponent -= 1
		smantissa = math.Round(mantissa * MDER_FLOAT_PRECISION)
		rmantissa = math.Round(mantissa) * MDER_FLOAT_PRECISION
		mdiff = math.Abs(smantissa - rmantissa)
	}

	binary.LittleEndian.PutUint32(result[1:], uint32(math.Round(float64(sign)*mantissa)))
	result[4] = byte(exponent)
	return result
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float32(math.Round(float64(float)*100) / 100)
}

func Float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
