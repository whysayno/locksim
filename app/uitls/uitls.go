package uitls

import (
	"bytes"
	"encoding/binary"
)

// 把一个整形数据转化成字节
func ConverNumberToBytes(v interface{}) []byte {
	btsByts := bytes.NewBuffer([]byte{})
	switch v.(type) {
	case int:
		binary.Write(btsByts, binary.BigEndian, v.(int))
	case int8:
		binary.Write(btsByts, binary.BigEndian, v.(int8))
	case int16:
		binary.Write(btsByts, binary.BigEndian, v.(int16))
	case int32:
		binary.Write(btsByts, binary.BigEndian, v.(int32))
	case uint:
		binary.Write(btsByts, binary.BigEndian, v.(uint))
	case uint16:
		binary.Write(btsByts, binary.BigEndian, v.(uint16))
	case uint8:
		binary.Write(btsByts, binary.BigEndian, v.(uint8))
	case uint32:
		binary.Write(btsByts, binary.BigEndian, v.(uint32))
	case uint64:
		binary.Write(btsByts, binary.BigEndian, v.(uint64))
	default:
		return nil
	}
	return btsByts.Bytes()
}

// CRC效验
func CRCCheck(buf []byte, length int) uint16 {
	var (
		CRC        uint16
		R          byte
		i, j, k, m int
	)
	if buf == nil || len(buf) < length || length <= 0 {
		return CRC
	}
	for i = 0; i < length; i++ {
		R = buf[i]
		for j = 0; j < 8; j++ {
			if R > 127 {
				k = 1
			} else {
				k = 0
			}
			R = R << 1
			if CRC > 0x7fff {
				m = 1
			} else {
				m = 0
			}
			if k+m == 1 {
				k = 1
			} else {
				k = 0
			}
			CRC = CRC << 1
			if k == 1 {
				CRC = CRC ^ 0x1021
			}
		}
	}
	return CRC
}
