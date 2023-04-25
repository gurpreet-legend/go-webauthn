package utils

import (
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
