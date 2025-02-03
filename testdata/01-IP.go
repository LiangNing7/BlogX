package main

import (
	"fmt"

	"github.com/LiangNing7/BlogX/core"
)

func main() {
	core.InitIPDB()
	fmt.Println(core.GetIpAddr("180.163.29.219"))
	fmt.Println(core.GetIpAddr("27.20.195.90"))
	fmt.Println(core.GetIpAddr("223.104.194.176"))
	fmt.Println(core.GetIpAddr("59.125.68.207"))
	fmt.Println(core.GetIpAddr("115.54.61.158"))
	fmt.Println(core.GetIpAddr("113.200.174.40"))
	fmt.Println(core.GetIpAddr("210.0.158.253"))
	fmt.Println(core.GetIpAddr("113.90.151.170"))
	fmt.Println(core.GetIpAddr("117.69.22.89"))
	fmt.Println(core.GetIpAddr("42.120.73.52"))
	fmt.Println(core.GetIpAddr("36.112.3.152"))
	fmt.Println(core.GetIpAddr("138.75.141.151"))
	fmt.Println(core.GetIpAddr("140.238.166.116"))
	fmt.Println(core.GetIpAddr("120.235.59.192"))
	fmt.Println(core.GetIpAddr("127.0.0.1"))
	fmt.Println(core.GetIpAddr("172.16.22.1"))
	fmt.Println(core.GetIpAddr("172.15.22.1"))
	fmt.Println(core.GetIpAddr("223.104.194.176"))
	fmt.Println(core.GetIpAddr("2231.104.194.176"))
	fmt.Println(core.GetIpAddr(""))
	fmt.Println(core.GetIpAddr("fe80::4464:1c7f:c207:6988%15"))
}
