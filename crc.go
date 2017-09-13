package main

import (
	"bytes"
)

func crc_16(data []byte) []byte {
	var buf bytes.Buffer
	i := 0
	j := 0
	var crc uint16 = 0xffff
	for i = 0; i < len(data); i++ {
		crc = crc ^ uint16(data[i])
		for j = 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ 0x8408
			} else {

				crc = (crc >> 1)
			}
		}

		//		fmt.Println(crc)
	}
	crc2 := crc >> 8
	last := crc << 8
	crc1 := last >> 8
	buf.WriteByte(byte(crc1))
	buf.WriteByte(byte(crc2))
	crc_byte := buf.Bytes()
	return crc_byte

}
