package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func writePkg(conn net.Conn, data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail! err=", err)
		return
	}
	//发送data本身
	n, err = conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Write fail! err=", err)
		return
	}
	return
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



