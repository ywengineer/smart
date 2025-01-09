package utility

import "encoding/binary"

func Int64ToBytes(odr binary.ByteOrder, i int64) []byte {
	var buf = make([]byte, 8)
	odr.PutUint64(buf, uint64(i))
	return buf
}

func Int32ToBytes(odr binary.ByteOrder, i int32) []byte {
	var buf = make([]byte, 4)
	odr.PutUint32(buf, uint32(i))
	return buf
}

func Int16ToBytes(odr binary.ByteOrder, i int16) []byte {
	var buf = make([]byte, 2)
	odr.PutUint16(buf, uint16(i))
	return buf
}

func BytesToInt64(odr binary.ByteOrder, buf []byte) int64 {
	return int64(odr.Uint64(buf))
}
