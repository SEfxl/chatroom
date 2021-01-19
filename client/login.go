package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"chatroom/common/message"
)

//写一个函数，完成登录
func Login(userId int, userPwd string) (err error) {
	//fmt.Printf("userId=%d, userPwd=%s\n", userId, userPwd)
	//return nil

	//1、链接到服务器端
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	//延时关闭
	defer conn.Close()

	//2、准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.LoginMesType

	//3、创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4、将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//5、将data赋给mes.data字段
	mes.Data = string(data)

	//6、将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7、data就是我们要发送的消息
	//7.1 先发送data的长度给服务器：先获取到data的长度-->然后转成一个表示长度的byte切片
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
	//fmt.Printf("客户端，发送消息的长度OK 长度=%d, 内容=%s", len(data), string(data))

	//7.2 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail! err=", err)
		return
	}

	//休眠20秒
	time.Sleep(time.Second * 20)
	fmt.Println("休眠20秒")
	//这里还需要处理服务器端返回的消息

	return
}
