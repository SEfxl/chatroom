package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

//定义一个全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle,maxActive int, idleTimeOut time.Duration)  {
	pool = &redis.Pool{
		MaxIdle:         maxIdle, //最大空闲连接数
		MaxActive:       maxActive, //表示和数据库的最大连接数，0表示没有限制
		IdleTimeout:     idleTimeOut, //最大空闲时间
		Dial: func() (conn redis.Conn, err error) { //初始化链接代码
			return redis.Dial("tcp", address)
		},
	}
}