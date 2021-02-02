package process2

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

//转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的onlineUsers map
	//将消息转发出去

	//取出mes的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes) //取出其中的消息,用来过滤掉自己
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(mes) //mes序列化之后直接做转发,不做任何处理
	if err != nil {
		fmt.Println("json.marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		//这里需要过滤掉自己,既不要再重新发给自己
		if id == smsMes.UserId {
			continue
		}

		this.SendMesToEachOnlineUser(data, up.Conn)
	}

}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个Transfer实例,发送data
	tf := utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=",err)
	}

}
