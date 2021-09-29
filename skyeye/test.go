package skyeye

import "github.com/youxiajinglin/skyeye-sdk-go/protobuf"

func ChatTest() *protobuf.ChatV3 {
	return &protobuf.ChatV3{
		Channel: "2",
		Content: "喴唁",
		Ip: "127.0.0.1",
		From: playerTest(),
	}
}

func playerTest() *protobuf.ChatUserV3 {
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

