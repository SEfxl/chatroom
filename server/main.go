package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"chatroom/common/message"
)

func main() {

	//提示信息
	fmt.Println("服务端在8889端口开始监听......")

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

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 4096)
	fmt.Println("读取客户端发送的数据....")
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//fmt.Println("读到的buf=", buf[0:4])
	//根据读到的长度,转换成uint32类型
	var pkgLen uint32 //表示实际读取到的内容的长度
	pkgLen = binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen]) //从conn套接字读取pkgLen的字节到buf中
	if uint32(n) != pkgLen || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	//把pkgLen,反序列化成-> message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return
}

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭
	defer conn.Close()

	//循环读取客户端发送的信息
	for {
		//读取包中的内容
		msg, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出,服务器也正常退出。。。")
				return
			} else {
				fmt.Println("readPkg err=", err)
				return
			}
		}
		fmt.Println("mes=", msg)
	}
}
