
```
    //初始化
    skyEye := skyeye.New([]byte{0xce, 0x35})
    //创建连接池，意外断开会自动会重连
    skyEye.NewPool(1, 2, func() (net.Conn, error) {
        return net.Dial("tcp", "192.168.1.110:8741")
    })

    go skyEye.Start()

    //获取聊天数据,并放入管道
    go func() {
        chat := getChat()
        for i := 0; i < 5; i++ {
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
```