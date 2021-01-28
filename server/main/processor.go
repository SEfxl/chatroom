package main

import (
	"chatroom/common/message"
	"chatroom/server/process"
	"chatroom/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn

}

//根据客户端发送消息种类不同，决定调用那个函数来处理
func (this *Processor)serverProcessMes(mes *message.Message) (err error) {
	//看看是否能接收到客户端发送的群发消息
	fmt.Println("mes=",mes)

	switch mes.Type {
		case message.LoginMesType: //处理登录的逻辑
			//创建一个UserProcess实例
			up := &process2.UserProcess{
				Conn:this.Conn,
			}
			err = up.ServerProcessLogin(mes)
		case message.RegisterMesType:
			//处理注册的逻辑
			//创建一个UserProcess实例
			up := &process2.UserProcess{
				Conn:this.Conn,
			}
			err = up.ServerProcessRegister(mes)
		default:
			fmt.Println("消息类型不存在,无法处理....")
	}

	return

}

func (this *Processor) process2() (err error) {
	//循环读取客户端发送的信息
	for {

		//创建一个Transfer实例完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出,服务器也正常退出。。。")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}

		//fmt.Println("mes=", msg)
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}

