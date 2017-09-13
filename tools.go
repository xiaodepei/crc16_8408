package main

import (
	"bytes"
	//	"fmt"
	"strconv"
)

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func reload_string_into_byte(data string) []byte {
	var buf bytes.Buffer
	for a := 0; a < len(data)/2; a++ {
		sep := Substr(data, a*2, 2)
		sep_int, _ := strconv.ParseInt(sep, 16, 16)
		buf.WriteByte(byte(sep_int))
	}
	return buf.Bytes()

}

func gettempature_cmd(epc []byte) []byte {
	cmd_1 := []byte{0x17, 0x00, 0x86, 0x06}
	part_1 := [][]byte{cmd_1, epc}
	a := bytes.Join(part_1, []byte(""))
	cmd_2 := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
	cmd := [][]byte{a, cmd_2}

	crc := crc_16(bytes.Join(cmd, []byte("")))
	all_cmd := [][]byte{bytes.Join(cmd, []byte("")), crc}

	result := bytes.Join(all_cmd, []byte(""))
	return result

}

func read_user_cmd(epc []byte) []byte {

	cmd_1 := []byte{0x18, 0x00, 0x02, 0x06}

	cmd_2 := []byte{0x03, 0xef, 0x01, 0x00, 0x00, 0x00, 0x00}

	part_1 := [][]byte{cmd_1, epc}

	a := bytes.Join(part_1, []byte(""))
	cmd := [][]byte{a, cmd_2}
	crc := crc_16(bytes.Join(cmd, []byte("")))
	all_cmd := [][]byte{bytes.Join(cmd, []byte("")), crc}
	result := bytes.Join(all_cmd, []byte(""))

	return result

}
