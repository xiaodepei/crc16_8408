package main

import (
	"fmt"
)

func main() {
	fmt.Println("Em4325芯片温度校准")
	fmt.Println("请输入温度数值（最小单位0.25度）：")
	fmt.Scanln(&stand_temp)
	read_tag()
}
