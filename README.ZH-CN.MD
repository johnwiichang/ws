# WS
WS 是一个扩展 Go 语言官方 WebSocket 工具库的扩展组件。目的旨协助开发者使用更熟悉的操作方式（HTTP 请求）来完成此类项目。

---

## 使用 HTTP 句柄风格的 WebSocket 处理器

如果你使用 `http.ListenAndServe` 来处理你的 Web 业务逻辑，你就可以通过这样来使用 WebSocket 赋予你程序全新的能量：

```
import "github.com/johnwiichang/ws"
//...
s = ws.NewWebSocketAdapter(ws.Adapter(func(conn *ws.Connection, body []byte) {
	fmt.Println(string(body))
	conn.Send([]byte("欢迎来到中国。"))
}))
http.ListenAndServe("0.0.0.0:9527", s)
```

你的业务逻辑将会被写成 `ws.Adapter`，其签名为 `Handler func(*ws.Connection, []byte)`. 若客户端发送 `你好中国` 至吾侪之服务器，你可以在终端中阅览到。使用这样的方式有效规避了帧处理以及避免了与客户端通信的额外转换：

```
conn.Send([]byte("欢迎来到中国。"))
```

客户端将会收到 `欢迎来到中国。`。

> 提示：
> 如果你觉得 `[]byte("欢迎来到中国。")` 啰嗦得不行，你可以直接使用 `conn.Send("欢迎来到中国。")` 来送出消息。本组件可以识别并转换它们到 `[]byte`。此外，如果你发送了一个结构体，那么将会转换为 JSON 格式并发送。但是你需要注意的是，如果你 **需要将字符串预先处理为 JSON 格式，你需要在发送前自行处理**。

## 自动序列化与反序列化
大多数场景中，我们使用结构化的传输方式来通信，我们也希望一个套接字能够承载更多的业务逻辑，彼此之间使用一个标记符号来甄别。用这个组件来实现这些愿望吧，代码应该是很开心的一件事才对。

天才第一步，创建一个实现了 `ws.Request` 接口的结构：
- WebSocket 通讯在建立起后无法通过路径来甄别操作名，所以你需要告诉 WS 你的业务名。
- 考虑有些时候，真正有意义的字段放置在其他的节点中。所以你可以告诉 WS 你的业务数据究竟在哪里。不过，你可以返回它自身来忽略这个你可能觉得画蛇添足的部分。

```
type (
	Request struct {
		OpToken string `json:"API"`
		Content struct {
			Server   string
			Location string
		} `json:"Body"`
	}
)

func (r Request) Action() string {
	return r.OpToken
}

func (r Request) Body() interface{} {
	return r.Content
}
```

之后，你就可以创建一个特别的 WebSocket 监听器：

```
s = ws.NewWebSocketService(Request{})
```

别担心，编译器会报告结构体是否是有效的。你可以将由 `Action() string` 说明的业务逻辑的标志符号注册到节点中。

例如：
```
s.RegisterEndpoint("/api/test", ws.Handler(func(conn *ws.Connection, body interface{}) {
	fmt.Println(body)
	conn.Send("テスト")
}))
```
这段片段将会收到一个请求实体，然后将客户端发送的消息的主体（Body）部分打印到屏幕上，最后告诉客户端 `テスト`。

> ⚠️ 警告:
> 只有 通过`NewWebSocketService` 创建的 `ws.Service` 才可以使用 `RegisterEndpoint` 方法，否则你将会一直遇到 panic 直到你放弃本扩展。

## 中途加密支持
在传统 HTTP API 服务中，你的数据需要手动加解密，繁琐不说，还有一大堆差不多的方法。为了减少这样的情况，也为了支持中途加密，使用本扩展时允许在通信时更新安全策略。意味着，天然支持非对等密钥交换的加密策略。

你只需要通过 `UpdateCodec` 就可以更新序列化反序列化器以及加密用服务。
不过，当心 `ws.Service` 以及 `ws.Connection` 都有这个方法。`Service` 的允许你指定新的默认的通信用编码译码器，另一个使用于通信中更换通信的策略的。

不妨试试利用这个特性在 WebSocket 中实现基于私有 PKI 设施的 ECDH 加密策略呗。