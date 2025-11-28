package util

import (
	"math/big"
	"strconv"
)

func getValAndBase(v string) (string, int) {
	if len(v) > 1 && (v[:2] == "0x" || v[:2] == "0X") {
		return v[2:], 16
	}
	return v, 10
}

func ToUint32(v string) uint32 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseUint(x, base, 32)
	if err != nil {
		return 0
	}
	return uint32(i)
}

func ToInt32(v string) int32 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseInt(x, base, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

func ToUint64(v string) uint64 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseUint(x, base, 64)
	if err != nil {
		return 0
	}
	return i
}

func ToUint8(v string) uint8 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseUint(x, base, 8)
	if err != nil {
		return 0
	}
	return uint8(i)
}

func ToInt64(v string) int64 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseInt(x, base, 64)
	if err != nil {
		return 0
	}
	return i
}

func ToBigInt(v string) *big.Int {
	x, base := getValAndBase(v)
	b := new(big.Int)
	b.SetString(x, base)
	return b
}

// BigIntToUintBytes 将大整数转换为指定长度的无符号字节切片
func BigIntToUintBytes(value *big.Int, byteLength int) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	// 获取大整数的字节表示
	bytes := value.Bytes()

	// 如果字节长度不足，在前面补0
	if len(bytes) < byteLength {
		padded := make([]byte, byteLength)
		copy(padded[byteLength-len(bytes):], bytes)
		return padded, nil
	}
	return bytes, nil
}

// Reverse 反转字节切片

func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
