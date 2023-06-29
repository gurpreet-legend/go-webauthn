package utils

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"time"
)

func GenerateUserID() uint64 {
	rand.Seed(time.Now().UnixNano())

	// Generate a random uint64 within the range of 0 to 9223372036854775807
	min := uint64(0)
	max := uint64(9223372036854775807)

	// Calculate the range of values
	rangeSize := max - min + 1

	// Generate a random value within the range
	userID := min + rand.Uint64()%rangeSize

	return userID
}

func RandomUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

func ConvertIntToByteArray(id uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(id))
	return buf
}

func ConvertByteArrayToInt(b []byte) uint64 {
	num, _ := binary.Uvarint(b)
	return num
}

func Base64ToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func StringToBase64(s string) []byte {
	res := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(res, []byte(s))
	return res
}
