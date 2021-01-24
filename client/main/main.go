package main

import (
	"chatroom/client/process"
	"fmt"
	"os"
)

//用户的ID和密码
var userId int
var userPwd string
var userName string

func main()  {
	//接受用户的选择
	var key int

	for true {
		fmt.Println("-------------欢迎登录多人聊天系统-------------")
		fmt.Println("\t\t\t 1 登录聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3)")

		fmt.Scanf("%d\n",&key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的ID")
			fmt.Scanf("%d\n",&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n",&userPwd)
			//完成登录
			//1、创建一个userProcess的实例
			up := &process.UserProcess{}
			up.Login(userId,userPwd)

			//loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户的ID")
			fmt.Scanf("%d\n",&userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n",&userPwd)
			fmt.Println("请输入用户的名字(nickname):")
			fmt.Scanf("%s\n",&userName)
			//2、调用userProcess完成注册的请求
			up := &process.UserProcess{}
			up.Register(userId,userPwd,userName)
			//loop = false
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("你的输入有误,请重新输入")
		}
	}

	//根据用户的输入,显示新的提示信息
	//if key == 1 {
	//	fmt.Println("请输入用户的ID")
	//	fmt.Scanf("%d\n",&userId)
	//	fmt.Println("请输入用户的密码")
	//	fmt.Scanf("%s\n",&userPwd)
	//	//因为使用了新的结构
	//
	//	//Login(userId,userPwd)
	//	//if err != nil {
	//	//	fmt.Println("登录失败")
	//	//}else {
	//	//	fmt.Println("登录成功")
	//	//}
	//}else if key == 2{
	//	fmt.Println("进行用户注册")
	//}


}
