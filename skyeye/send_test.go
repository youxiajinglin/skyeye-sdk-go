package skyeye

import (
	"fmt"
	"github.com/youxiajinglin/skyeye-sdk-go/skyeye/protobuf"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSkyEye_Send(t *testing.T) {
	//初始化
	skyEye := New([]byte{0xce, 0x35})
	//创建连接池，意外断开会自动会重连
	skyEye.NewPool(1, 2, func() (net.Conn, error) {
		return net.Dial("tcp", "192.168.1.110:8741")
	})

	go skyEye.Start()

	wg := sync.WaitGroup{}
	//回执处理
	go skyEye.Reply(func(v3 *protobuf.ChatV3) {
		wg.Done()
		fmt.Printf("reply ID：%s 审核结果: %d \r", v3.GetId() ,v3.GetStatus())
	})

	//获取聊天数据,并放入管道
	chat := getChat()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		n := strconv.Itoa(i)
		chat.Id = n
		skyEye.ChatToChan(chat)
		time.Sleep(2*time.Second)
	}

	wg.Wait()
}

func getChat() *protobuf.ChatV3 {
	return &protobuf.ChatV3{
		Channel: "2",
		Content: "l1109528",
		Ip: "127.0.0.1",
		From: getPlayer(),
	}
}

func getPlayer() *protobuf.ChatUserV3 {
	return &protobuf.ChatUserV3{
		UserId: "5566",
		PlayerId: "7788",
		Nickname: "测试用户",
		ZoneId: "1",
		ZoneName: "",
		ServerId: "test",
		Level: 100,
		VipLevel: 5,
		Extra: "",
	}
}

