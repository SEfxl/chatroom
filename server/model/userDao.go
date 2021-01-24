package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//在服务启动后就初始化一个userDao实例
//把它做成全局的变量,在需要和redis操作时就直接使用即可
var (
	MyUserDao *UserDao
)


//userDao结构体完成对user结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式,创建一个userDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {

	userDao = &UserDao{
		pool:pool,
	}

	return
}


//根据用户ID 返回一个user实例
func (this *UserDao)getUserById(conn redis.Conn, id int) (user *message.User, err error) {

	res , err := redis.String(conn.Do("HGet","users",id))
	if err != nil {
		if err == redis.ErrNil { //表示在users哈希中,没有找到对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &message.User{}
	//这里需要反序列化res成一个User实例
	err = json.Unmarshal([]byte(res),user)
	if err != nil {
		fmt.Println("json.Unmarshal err=",err)
		return
	}

	return
}

//完成登录的校验
//1、Login完成对用户的校验
//2、如果用户的id和pwd都正确,则返回一个user实例
//3、如果用户的id和pwd有错误,则返回对应的错误信息
func (this *UserDao)Login(userId int, userPwd string) (user *message.User, err error)  {

	//先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	user,err = this.getUserById(conn,userId)
	if err != nil {
		return
	}

	//证明用户获取到
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}


//用户注册
func (this *UserDao)Register(user *message.User) (err error)  {

	//先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	_,err = this.getUserById(conn,user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	//这时说明id在redis中不存在,则可以完成注册
	data , err := json.Marshal(user) //结构体序列化为字符串
	if err != nil {
		return
	}

	//入库
	_, err = conn.Do("HSet","users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=",err)
		return
	}

	return
}













