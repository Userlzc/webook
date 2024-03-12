package main

import (
	wire2 "project/wire"
)

// 双击shift可查找任何东西
// ctrl +F 查找该文件中匹配的东西

func main() {
	server := wire2.InitWebServer()
	err := server.Run(":8081")
	if err != nil {
		panic("端口可能被占用")
	}

}
