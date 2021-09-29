
```
    //初始化
    skyEye := skyeye.New([]byte{0xce, 0x35})
    //创建连接池，意外断开会自动会重连
    skyEye.NewPool(1, 1, func() (net.Conn, error) {
        //设置连接超时2秒
        return net.DialTimeout("tcp", "192.168.1.110:8741", time.Second*2)
    })

    go skyEye.Start()

    //回执处理
    go skyEye.Reply(func(resp *skyeye.Response) {
        if resp.Err != nil {
            fmt.Println(resp.Err)
            return
        }
        fmt.Printf("reply ID：%s 审核结果: %d \r", resp.Data.GetId() ,resp.Data.GetStatus())
    })
    
    //获取聊天数据,并放入管道
    chat := skyeye.ChatTest()
    for i := 0; i < 100; i++ {
        n := strconv.Itoa(i)
        chat.Id = n
        skyEye.Push(chat)
        time.Sleep(1*time.Second)
    }
    
    select {}
```