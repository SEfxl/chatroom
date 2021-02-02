package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

func main() {

	//1、服务器启动时,初始化redis的连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	//2、初始化userDao
	initUserDao()

	//提示信息
	fmt.Println("服务器【新的结构】在8889端口开始监听......")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	//一旦监听成功，就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器......")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}
}

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭
	defer conn.Close()

	//这里调用总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯的协程出错", err)
		return
	}

}

//写一个函数完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool本身就是一个全局的变量
	model.MyUserDao = model.NewUserDao(pool)
}
