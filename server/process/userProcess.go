package process2

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserId int
}

//编写通知所有在线用户的方法
//userId 通知其他在线用户,我上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {

	//遍历onlineUsers 然后一个一个发送NotifyUserStatusMes消息
	for id, up := range userMgr.onlineUsers {
		if id == userId { //过滤自己
			continue
		}

		//开始通知,NotifyMeOnline的this就是up
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装消息NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)

	//对mes再次序列化,准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//发送,创建一个Transfer实例发送
	tf := utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

//用户注册的请求
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1、先从mes中取出mes.Data 并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//1、先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	//2、再声明一个LoginResMes,并完成赋值
	var registerResMes message.RegisterResMes

	//需要到redis去完成注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发送未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	//3、将loginResMes序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4、将data赋值给resMes
	resMes.Data = string(data)

	//5、对resMes进行序列化 准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//6、发送data 将其分装到writePkg中
	//因为使用了分层的模式,我们先创建一个Transfer实例,然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return
}

//专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//1、先从mes中取出mes.Data 并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	//1、先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//2、再声明一个LoginResMes,并完成赋值
	var loginResMes message.LoginResMes

	//需要到redis去验证用户
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误...."
		}

	} else {
		loginResMes.Code = 200

		//这里,因为用户已经登录成功,我们就把该登录成功的用户放入到userMgr中
		//将登录成功的用户的userId赋给 this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)

		//登陆成功将自己登陆成功的消息推送给其他在线用户
		this.NotifyOthersOnlineUser(loginMes.UserId)

		//将当前在线用户id,放入到loginResMes.UsersId
		//遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}

		fmt.Println(user.UserName, "登录成功")
	}

	//如果用户id=100 密码=123456 认为合法，否则不合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//	loginResMes.Code = 200
	//}else{
	//	loginResMes.Code = 500 //500表示该用户不存在
	//	loginResMes.Error = "该用户不存在,请注册再使用..."
	//}

	//3、将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4、将data赋值给resMes
	resMes.Data = string(data)

	//5、对resMes进行序列化 准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//6、发送data 将其分装到writePkg中
	//因为使用了分层的模式,我们先创建一个Transfer实例,然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
