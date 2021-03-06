package message

//确定一些消息类型
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

//定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	userBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息的类型
	Data string `json:"data"` //消息的数据
}

//定义具体的消息-登陆
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户ID
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code    int    `json:"code"`  //返回状态码 500表示该用户未注册 200表示登陆成功
	Error   string `json:"error"` //返回错误信息
	UsersId []int  //增加字段,保存用户id的切片
}

//定义具体的消息-注册
type RegisterMes struct {
	User User `json:"user"` //类型就是User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码 400表示该用户已经占用 200表示注册成功
	Error string `json:"error"` //返回错误信息
}

//定义具体的消息-服务端通知客户端
//配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户状态
}



//增加一个SmsMes(发送消息)
type SmsMes struct {
	Content string `json:"content"` //内容
	User //匿名结构体,继承  User User
 }










