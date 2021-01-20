package message

//确定一些消息类型

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
)

type Message struct {
	Type string `json:"type"` //消息的类型
	Data string `json:"data"` //消息的数据
}

//定义具体的消息
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户ID
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code  int    `json:"code"`  //返回状态码 500表示该用户未注册 200表示登陆成功
	Error string `json:"error"` //返回错误信息
}

type RegisterMes struct {

}