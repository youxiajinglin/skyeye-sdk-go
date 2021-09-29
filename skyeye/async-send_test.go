package skyeye

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"
)

func TestSkyEye_Send(t *testing.T) {
	//初始化
	skyEye := New([]byte{0xce, 0x35})
	//创建连接池，意外断开会自动会重连
	skyEye.NewPool(1, 2, func() (net.Conn, error) {
		//设置连接超时2秒
		conn, err := net.Dial("tcp", "192.168.1.110:8741")
		if err != nil{
			panic(err)
		}
		//err = conn.SetWriteDeadline(time.Now().Add(1*time.Second))
		return conn, err
	})

	go skyEye.Start()

	wg := sync.WaitGroup{}
	//TCP的回执需根据id来异步处理，因为当前的回执可能不是当前这一条的
	skyEye.Reply(func(resp *Response) {
		wg.Done()
		if resp.Err != nil {
			fmt.Println(resp.Err)
			return
		}
		fmt.Printf("reply ID：%s 审核结果: %d \r", resp.Data.GetId() ,resp.Data.GetStatus())
	})

	//获取聊天数据,并放入管道
	chat := ChatTest()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		n := strconv.Itoa(i)
		chat.Id = n
		skyEye.Push(chat)
		//time.Sleep(1*time.Second)
	}

	wg.Wait()
}