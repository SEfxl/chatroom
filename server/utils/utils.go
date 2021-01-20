package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	Conn  net.Conn    //连接
	Buf   [8192]byte  //传输时使用的缓冲
}

//读取包中的内容 读取数据包分装成一个函数,返回message err
func (this *Transfer)ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 4096)
	fmt.Println("读取客户端发送的数据....")
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//fmt.Println("读到的buf=", buf[0:4])
	//根据读到的长度,转换成uint32类型
	var pkgLen uint32 //表示实际读取到的内容的长度
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen]) //从conn套接字读取pkgLen的字节到buf中
	if uint32(n) != pkgLen || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	//把pkgLen,反序列化成-> message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return
}


func (this *Transfer)WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail! err=", err)
		return
	}
	//发送data本身
	n, err = this.Conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Write fail! err=", err)
		return
	}
	return
}