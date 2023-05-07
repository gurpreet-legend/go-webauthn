package utils

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
)

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

func Base64ToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func StringToBase64(s string) []byte {
	res := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(res, []byte(s))
	return res
}
