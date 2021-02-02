package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

//这个mes一定是SmsMes
func outputGroupMes(mes *message.Message)  {
	//显示即可

	//1、反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data),&smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=",err)
		return
	}

	//2、显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s",smsMes.UserId,smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
