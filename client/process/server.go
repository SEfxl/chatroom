package process

import (
	"encoding/json"

	"chatroom/client/utils"
	"chatroom/common/message"

	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面。。。
func ShowMenu() {
	fmt.Println("--------恭喜xxx登录成功--------")
	fmt.Println("--------1、显示在线用户列表--------")
	fmt.Println("--------2、发送消息--------")
	fmt.Println("--------3、信息列表--------")
	fmt.Println("--------4、退出系统--------")
	fmt.Println("请选择(1-4):")

	var key int
	var content string

	//因为总会使用到smsProcess实例,因此将其定义在switch外部
	smsProcess := &SmsProcess{}

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("你想对大家说点什么:)")
		fmt.Scanf("%s\n",&content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择退出了系统.....")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不对..")
	}
}

//和服务器端保持通讯
func serverProcessMes(conn net.Conn) {

	//创建一个transfer实例,不停的读取服务器发送的消息
	//只要对方不关闭连接,就会一直读取,并且对程序的效率没有影响
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		//如果读取到消息,又是下一步处理逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线
			//1、取出NotifyUserStatusMesType
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2、把这个用户的信息,状态保存到客户端map[int]user中
			updateUserStatus(&notifyUserStatusMes)
		default:
			fmt.Println("服务器端返回了一个未知的消息类型")
		}

		//fmt.Printf("read from server! mes=%v\n",mes)

	}
}
