package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"time"
	//	"fmt"
	"strconv"

	"github.com/tarm/goserial"
	"gopkg.in/redis.v5"
)

var ser = &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
var s, _ = serial.OpenPort(ser)

var stand_temp float64

var Client13 = redis.NewClient(&redis.Options{
	Addr:     "192.168.1.137:6379",
	Password: "xiaodepe", //  password set
	DB:       0,          // use default DB
})

//var cha1 = make(chan string, 10)

func read_tag() {
	for {

		//读写器查询标签指令（快速）
		data := []byte{0x06, 0x00, 0x01, 0x01, 0xff, 0x6c, 0x47}
		//读写器查询标签指令（慢速）
		//		data := []byte{0x06, 0x00, 0x01, 0x04, 0xff, 0xd4, 0x39}
		s.Write(data)
		time.Sleep(800000000)
		buf := make([]byte, 2048)
		n, _ := s.Read(buf)
		p := hex.EncodeToString(buf)
		len_str := n*2 - 16

		data1 := Substr(p, 0, n*2)

		var temp string
		for a := 0; a < len_str/44; a++ {

			tag_1 := Substr(data1, 44*a+14, 24)

			//			tag_1_ant := Substr(data1, 44*a+8, 2)
			//			fmt.Println(tag_1)
			if temp != tag_1 {
				temp = tag_1

				//				fmt.Println(tag_1)
				byte_epc := reload_string_into_byte(tag_1)
				if byte_epc[1] == 0x80 {

					if Client13.Exists(tag_1).Val() == false {

						get_temp := gettempature_cmd(byte_epc)
						gettempature(get_temp, tag_1)
						fmt.Println("原温度如上")
						rewrite_data := get_rewrite_user(byte_epc)
						fmt.Println("得到重写数据")
						write_user_cmd := write_user_zero_cmd(byte_epc, rewrite_data)
						fmt.Println(rewrite_data)
						fmt.Println("得到重写指令")
						write_user_zero(write_user_cmd)
						clean_temp := gettempature(get_temp, tag_1)
						fmt.Println("得到清除原校准的温度")
						fmt.Println(clean_temp)
						rewrite_data = jiaozhun(clean_temp, rewrite_data)
						write_user_cmd = write_user_zero_cmd(byte_epc, rewrite_data)
						fmt.Println(rewrite_data)
						write_user_zero(write_user_cmd)
						fmt.Println("已校准完毕")
						aa := gettempature(get_temp, tag_1)

						fmt.Printf("校准后温度为：%s,", aa)
						Client13.Set(tag_1, aa, 0)

						//						Client13.HSet(tag_1, "校准前", string(bb))nt13.HSet(tag_1, "校准后", string(aa))

					}

				}

			}
		}
	}
}

func gettempature(cmd []byte, epc string) float64 {
	s.Write(cmd)
	var temp_real float64
	time.Sleep(800000000)
	buf := make([]byte, 1024)
	s.Read(buf)
	if buf[0] == 0x0d {
		p := hex.EncodeToString(buf)
		temp := Substr(p, 10, 2)
		temp_int, _ := strconv.ParseInt(temp, 16, 10)
		temp_real = float64(temp_int) * 0.25
		fmt.Println(epc)
		fmt.Println(temp_real)

	}
	return temp_real
}

func get_rewrite_user(epc []byte) []byte {
	var buff bytes.Buffer
	var rewrite_byte []byte
	cmd_read := read_user_cmd(epc)
	s.Write(cmd_read)
	time.Sleep(800000000)
	buf := make([]byte, 1024)
	s.Read(buf)
	if buf[0] == 0x07 {
		p := hex.EncodeToString(buf)
		user_keep1 := Substr(p, 8, 2)
		user_keep1_int, _ := strconv.ParseInt(user_keep1, 16, 16)
		user_rewrite := Substr(p, 10, 2)
		user_rewrite_int, _ := strconv.ParseInt(user_rewrite, 16, 16)
		user_rewrite_int_keep := user_rewrite_int >> 5
		user_rewrite_int_keep = user_rewrite_int_keep << 5

		buff.WriteByte(byte(user_keep1_int))
		buff.WriteByte(byte(user_rewrite_int_keep))
		rewrite_byte = buff.Bytes()

	}
	return rewrite_byte

}

func write_user_zero_cmd(epc []byte, rewrite_data []byte) []byte {
	cmd_1 := []byte{0x1a, 0x00, 0x03, 0x01, 0x06}
	cmd_2 := []byte{0x03, 0xef}
	cmd_3 := []byte{0x00, 0x00, 0x00, 0x00}
	part_1 := [][]byte{cmd_1, epc}
	a := bytes.Join(part_1, []byte(""))
	cmd := [][]byte{a, cmd_2}
	cmd_13 := bytes.Join(cmd, []byte(""))
	cmd_4 := [][]byte{cmd_13, rewrite_data}
	cmd_5 := bytes.Join(cmd_4, []byte(""))
	cmd_final := [][]byte{cmd_5, cmd_3}
	crc := crc_16(bytes.Join(cmd_final, []byte("")))
	all_cmd := [][]byte{bytes.Join(cmd_final, []byte("")), crc}
	result := bytes.Join(all_cmd, []byte(""))
	return result

}

func write_user_zero(cmd []byte) {
	s.Write(cmd)
	time.Sleep(800000000)
	buf := make([]byte, 1024)
	s.Read(buf)
	fmt.Println(buf[3])
	if buf[3] == 0x00 {
		fmt.Println("清除原校准数据")

	} else {
		return
	}

}

func jiaozhun(clean_temp float64, keep_rewrite_data []byte) []byte {
	check_temp := clean_temp - stand_temp
	if check_temp > 0 {

		flag := 1 << 4

		data := (math.Abs(check_temp)) / 0.25

		if data > 11 {

			data = 11
		}
		fmt.Println("dsfd")

		a := int(data) - 1 ^ 15

		tempatuer_part := flag ^ a
		keep_rewrite_data[1] = keep_rewrite_data[1] ^ byte(tempatuer_part)

	}

	if check_temp < 0 {

		data := (math.Abs(check_temp)) / 0.25

		keep_rewrite_data[1] = keep_rewrite_data[1] ^ byte(data)

	}
	return keep_rewrite_data

}
