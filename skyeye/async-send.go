package skyeye

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/pool"
	"github.com/golang/protobuf/proto"
	"github.com/youxiajinglin/skyeye-sdk-go/protobuf"
	"net"
	"time"
)

type SkyEye struct {
	id []byte
	list chan *protobuf.ChatV3
	result chan *Response
	pool pool.Pool
}

type Response struct {
	Data *protobuf.ChatV3
	Err error
}

//id 天眼后台生成的2 bytes的ID
func New(id []byte) *SkyEye {
	return &SkyEye{
		id: id,
		list: make(chan *protobuf.ChatV3, 100),
		result: make(chan *Response, 100),
	}
}

//initialCap 初始连接数
//maxCap 最大连接数
func (c *SkyEye) NewPool(initialCap, maxCap int, f func() (net.Conn, error)) {
	pool, err := pool.NewChannelPool(initialCap, maxCap, f)
	if err != nil {
		fmt.Errorf("create connect pool fail: %v", err)
	}
	c.pool = pool
}

func (c *SkyEye) Start() {
	for {
		select {
		case chat := <-c.list:

			if chat != nil {
				go func(chat *protobuf.ChatV3) {
					var (
						buf *protobuf.ChatV3
						err error
					)

					data, err := c.Send(chat)
					if err != nil {
						return
					}
					buf, err = c.byteToProtobuf(data)
					if err != nil {
						return
					}
					defer func() {
						c.result <-&Response{
							Data: buf,
							Err:  err,
						}
					}()
				}(chat)
			}
		}
	}
}

func (c *SkyEye) Send(chat *protobuf.ChatV3) ([]byte, error) {
	encode, err := proto.Marshal(chat)
	if err != nil {
		return nil, err
	}
	buffer := Encode(c.id, encode)
	conn, err := c.pool.Get()

	defer func() {
		//释放连接，放入连接池
		if conn != nil {
			conn.Close()
		}
	}()
	//连接可能为nil
	if err != nil || conn == nil {
		//重新放回管道
		c.Push(chat)
		return nil, errors.New("connect fail, rejoin chan")
	}
	err = conn.SetReadDeadline(time.Now().Add(1*time.Second))
	if err != nil {
		//重新放回管道
		c.Push(chat)
		return nil, errors.New("SetWriteDeadline fail")
	}
	_, err = conn.Write(buffer)
	//写入失败，可能连接已断开
	if err != nil {
		//重新放回管道
		c.Push(chat)
		return nil, err
	}
	reader := bufio.NewReader(conn)
	result, err := Decode(reader)
	//接收失败，代表连接可能断开，需重连
	if err != nil {
		//将失效的连接关闭
		//pc.MarkUnusable()
		//pc.Close()

		return nil, err
	}

	//回执加入chan
	return result, nil
}

func (c *SkyEye) Reply(call func(resp *Response)) {
	for {
		select {
		case result := <-c.result:
			call(result)
		}
	}
}

func (c *SkyEye) Push(chat *protobuf.ChatV3) {

	fmt.Println("push", chat.Id)

	c.list <-chat
}

func (c *SkyEye) byteToProtobuf(result []byte) (*protobuf.ChatV3, error) {
	chat := &protobuf.ChatV3{}
	err := proto.Unmarshal(result, chat)
	if err != nil {
		return nil, err
	}
	return chat, nil
}