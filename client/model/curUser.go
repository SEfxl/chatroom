package model

import (
	"chatroom/common/message"
	"net"
)

//设置该结构体就是维护了一个客户端的连接,为了和服务器通信
//因为在客户端,我们很多地方会使用到curUser 我们将其作为一个全局的
type CurUser struct {
	Conn net.Conn
	message.User
}


