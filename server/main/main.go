package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

func main() {

	//1、服务器启动时,初始化redis的连接池
	initPool("localhost:6379",16,0,300 * time.Second)
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
		Conn:conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯的协程出错",err)
		return
	}

}

//写一个函数完成对UserDao的初始化任务
func initUserDao()  {
	//这里的pool本身就是一个全局的变量
	model.MyUserDao = model.NewUserDao(pool)
}

//专门处理登录请求
//func serverProcessLogin(conn net.Conn, mes *message.Message) (err error)  {
//	//1、先从mes中取出mes.Data 并直接反序列化成LoginMes
//	var loginMes message.LoginMes
//	err = json.Unmarshal([]byte(mes.Data),&loginMes)
//	if err != nil {
//		fmt.Println("json.Unmarshal fail err=",err)
//		return
//	}
//
//	//1、先声明一个 resMes
//	var resMes message.Message
//	resMes.Type = message.LoginResMesType
//
//	//2、再声明一个LoginResMes
//	var loginResMes message.LoginResMes
//
//	//如果用户id=100 密码=123456 认为合法，否则不合法
//	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
//		loginResMes.Code = 200
//	}else{
//		loginResMes.Code = 500 //500表示该用户不存在
//		loginResMes.Error = "该用户不存在,请注册再使用..."
//	}
//
//	//3、将loginResMes序列化
//	data , err := json.Marshal(loginResMes)
//	if err != nil {
//		fmt.Println("json.Marshal fail",err)
//		return
//	}
//
//	//4、将data赋值给resMes
//	resMes.Data = string(data)
//
//	//5、对resMes进行序列化 准备发送
//	data , err = json.Marshal(resMes)
//	if err != nil {
//		fmt.Println("json.Marshal fail",err)
//		return
//	}
//
//	//6、发送data 将其分装到writePkg中
//	err = writePkg(conn,data)
//	return
//}

//func writePkg(conn net.Conn, data []byte) (err error) {
//	 //先发送一个长度给对方
//	var pkgLen uint32
//	pkgLen = uint32(len(data))
//	var buf [4]byte
//	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
//	//发送长度
//	n, err := conn.Write(buf[0:4])
//	if n != 4 || err != nil {
//		fmt.Println("conn.Write fail! err=", err)
//		return
//	}
//	//发送data本身
//	n, err = conn.Write(data)
//	if uint32(n) != pkgLen || err != nil {
//		fmt.Println("conn.Write fail! err=", err)
//		return
//	}
//	return
//}

//根据客户端发送消息种类不同，决定调用那个函数来处理
//func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
//	switch mes.Type {
//		case message.LoginMesType: //处理登录的逻辑
//			err = serverProcessLogin(conn,mes)
//		case message.RegisterMesType:
//			//处理注册的逻辑
//		default:
//			fmt.Println("消息类型不存在,无法处理....")
//	}
//
//	return
//}

//func readPkg(conn net.Conn) (mes message.Message, err error) {
//	buf := make([]byte, 4096)
//	fmt.Println("读取客户端发送的数据....")
//	_, err = conn.Read(buf[:4])
//	if err != nil {
//		//err = errors.New("read pkg header error")
//		return
//	}
//	//fmt.Println("读到的buf=", buf[0:4])
//	//根据读到的长度,转换成uint32类型
//	var pkgLen uint32 //表示实际读取到的内容的长度
//	pkgLen = binary.BigEndian.Uint32(buf[0:4])
//
//	//根据pkgLen读取消息内容
//	n, err := conn.Read(buf[:pkgLen]) //从conn套接字读取pkgLen的字节到buf中
//	if uint32(n) != pkgLen || err != nil {
//		//err = errors.New("read pkg body error")
//		return
//	}
//
//	//把pkgLen,反序列化成-> message.Message
//	err = json.Unmarshal(buf[:pkgLen], &mes)
//	if err != nil {
//		fmt.Println("json.Unmarshal err=", err)
//		return
//	}
//
//	return
//}


