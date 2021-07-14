package chat

import (
	"fmt"
	"github.com/dong568789/skyeye-sdk-go/chat/protobuf"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestSkyEye_Send(t *testing.T) {
	//初始化
	skyEye := NewSkyEye([]byte{0xce, 0x35})
	//创建连接池，意外断开会自动会重连
	skyEye.NewPool(1, 2, func() (net.Conn, error) {
		return net.Dial("tcp", "192.168.1.110:8741")
	})

	go skyEye.Start()

	//获取聊天数据,并放入管道
	go func() {
		chat := getChat()
		for i := 0; i < 10; i++ {
			n := strconv.Itoa(i)
			chat.Id = n
			skyEye.ChatToChan(chat)
			time.Sleep(2*time.Second)
		}
	}()

	//回执处理
	go skyEye.Reply(func(v3 *protobuf.ChatV3) {
		fmt.Printf("reply ID：%s 审核结果: %d \r", v3.GetId() ,v3.GetStatus())
	})
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
