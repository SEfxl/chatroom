package main

import "fmt"

//写一个函数，完成登录
func Login(userId int,userPwd string) (err error) {
	fmt.Printf("userId=%d, userPwd=%s\n",userId,userPwd)
	return nil
}