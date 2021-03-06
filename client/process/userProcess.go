package process

import (
	"chatroom/client/utils"
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//暂时不需要字段
type UserProcess struct {
}

func (this *UserProcess) Register(userId int, userPwd string,
	userName string) (err error) {

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
	mes.Type = message.RegisterMesType

	//3、创建一个LoginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4、将registerMes序列化
	data, err := json.Marshal(registerMes)
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

	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	mes, err = tf.ReadPkg() //mes就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//将mes的data部分反序列化成 LoginResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功,重新登录一下")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}

func (this *UserProcess) Login(userId int, userPwd string) (err error) {
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

	fmt.Printf("客户端，发送消息的长度OK 长度=%d, 内容=%s", len(data), string(data))

	//7.2 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail! err=", err)
		return
	}

	//休眠20秒
	//time.Sleep(time.Second * 20)
	//fmt.Println("休眠20秒")

	//这里还需要处理服务器端返回的消息
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg() //mes就是
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//将mes的data部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//fmt.Println("登录成功")

		//初始化curUser
		curUser.Conn = conn
		curUser.UserId = userId
		curUser.UserStatus = message.UserOnline


		//显示当前在线用户列表,遍历loginResMes.UsersId
		fmt.Println("当前在线用户列表如下")
		for _, v := range loginResMes.UsersId {

			//如果我们要求不显示自己在线,下面我们增加一个代码
			if v == userId {
				continue
			}

			fmt.Println("用户id:\t", v)

			//完成客户端onlineUsers的初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user

		}
		fmt.Println()

		//这里需要在客户端启动一个协程，该协程保持和服务器端的通讯,
		//如果服务器有数据推送给客户端,则接受并显示在客户端的终端
		go serverProcessMes(conn)

		//1、显示登录成功后的菜单[循环显示]
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}
