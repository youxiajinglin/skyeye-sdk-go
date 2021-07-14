package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dong568789/skyeye-sdk-go/chat/protobuf"
	"github.com/fatih/pool"
	"github.com/golang/protobuf/proto"
	"net"
)

type SkyEye struct {
	id []byte
	List chan *protobuf.ChatV3
	Result chan []byte
	pool pool.Pool
}

//id 天眼后台生成的2 bytes的ID
func NewSkyEye(id []byte) *SkyEye {
	return &SkyEye{
		id: id,
		List: make(chan *protobuf.ChatV3, 100),
		Result: make(chan []byte, 100),
	}
}

func (c *SkyEye) ChatToChan(chat *protobuf.ChatV3) {
	c.List <-chat
}

func (c *SkyEye) Start() {
	for {
		select {
		case chat := <-c.List:
			go c.send(chat)
		}
	}
	close(c.List)
}

//initialCap 初始连接数
//maxCap 最大连接数
func (c *SkyEye) NewPool(initialCap, maxCap int,f func() (net.Conn, error)) {
	pool, err := pool.NewChannelPool(initialCap, maxCap, f)
	if err != nil {
		fmt.Errorf("create connect pool fail: ", err)
	}
	c.pool = pool
}

func (c *SkyEye) Reply(call func(v3 *protobuf.ChatV3)) {
	for {
		select {
		case result := <-c.Result:
			chat := &protobuf.ChatV3{}
			proto.Unmarshal(result, chat)
			call(chat)
		}
	}
	close(c.Result)
}


func (c *SkyEye) send(chat *protobuf.ChatV3)  {
	fmt.Printf("send id:%s content:%s \n", chat.GetId(), chat.GetContent())
	encode, err := proto.Marshal(chat)
	if err != nil {
		fmt.Println("parse protobuf fail", err)
		return
	}
	buffer := c.encode(encode)
	conn, _ := c.pool.Get()
	//连接可能为nil
	if conn == nil {
		return
	}
	_, err = conn.Write(buffer)
	if err != nil {
		fmt.Errorf("send fail: %v", err)
		return
	}
	reader := bufio.NewReader(conn)
	result, err := c.decode(reader)
	//接收失败，代表连接可能断开，需重连
	if err != nil {
		//重新放回管道
		c.List <-chat
		//将失效的连接关闭
		if pc, ok := conn.(*pool.PoolConn); ok {
			pc.MarkUnusable()
			pc.Close()
		}
	}
	//释放连接，放加连接池
	conn.Close()

	//回执加入chan
	c.Result <-result
}

func (c *SkyEye) encode(content []byte) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, c.id)
	binary.Write(buffer, binary.BigEndian, uint32(len(content)))
	binary.Write(buffer, binary.BigEndian, content)
	return buffer.Bytes()
}

func (c *SkyEye) decode(reader *bufio.Reader) ([]byte, error) {
	haedByte, err := reader.Peek(6)
	if err != nil {
		return nil, err
	}
	lengthBuff := bytes.NewBuffer(haedByte[2:6])

	var lenght int32
	err = binary.Read(lengthBuff, binary.BigEndian, &lenght)
	if err != nil {
		return nil, err
	}
	if lenght + 6 > int32(reader.Buffered()) {
		return nil, err
	}

	pack := make([]byte, int(6 + lenght))
	_, err = reader.Read(pack)
	if err != nil {
		return nil, err
	}
	return pack[6:], nil
}