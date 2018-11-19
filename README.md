# WS
WS is an extension package for golang official WebSocket tool pack. It can help developers complete the project with familiar operations just like the http request one.

---

## HTTP Handler-like WebSocket Processor

If you use `http.ListenAndServe` to handle your logical web service, and you can just add the websocket one like:

```
import "github.com/johnwiichang/ws"
//...
s = ws.NewWebSocketAdapter(ws.Adapter(func(conn *ws.Connection, body []byte) {
	fmt.Println(string(body))
	conn.Send([]byte("欢迎来到中国。"))
}))
http.ListenAndServe("0.0.0.0:9527", s)
```

Your logic and be written into `ws.Adapter` which the signature is `Handler func(*ws.Connection, []byte)`. If a client is sending text `你好中国` to your server, you will find these characters in your terminal.
It’s pretty nice to ignore frame operation and response the client directly without any transforms:

```
conn.Send([]byte("欢迎来到中国。"))
```

The client will receive a string which is `欢迎来到中国。`

> P.S.:
> if you think that `[]byte("欢迎来到中国。")` is pretty verbose, you can use `conn.Send("欢迎来到中国。")` directly. WS can identify it and trans theme into `[]byte`. Also, if you send a struct, WS will transform it into JSON format, but you **should marshal string into JSON format if you’d like to send a JSON formatted string to client**.

## Serialise & Deserialise Automatically
But in most scenarios, we use structured transport structures to communicate. And we may hope multiple sets of business logic with instruction flags and be handled by one socket. Don’t worry, let’s make the code happier.
First of all, define a struct which implemented `ws.Request` interface.
- The WebSocket communication cannot identify the action name by request path after connecting. So, you need to tell WS, which field is the name of action.
- We considered that operation identifier might just a no meaning field, and your data might storage in an element. So you can tell WS, which field is most important and have to be calculated. Otherwise, you can just return itself to ignore.
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
After that, you can create a special listener using WebSocket:
```
s = ws.NewWebSocketService(Request{})
```
The compiler will check whether the strict is valid or not, and you can register your business logic into endpoint which described by `Action() string`.
Demo:
```
s.RegisterEndpoint("/api/test", ws.Handler(func(conn *ws.Connection, body interface{}) {
	fmt.Println(body)
	conn.Send("テスト")
}))
```
This snippet can receive a `Request` entity and give a string back after printing entity on the screen.

> ⚠️ Notice:
> The `RegisterEndpoint` is available for `ws.Service` created by `NewWebSocketService`. If you try to assign an endpoint to an instance created by `NewWebSocketAdapter`, a beautiful panic will flash to you immediately.

## Midway Encryption Support
In traditional http api service, your data should be encrypted / decrypted by coder-self. It’s verbose and invoke lots of identical functions. To reduce and support midway encryption security policy, you can update your security policy during the communication. It means, non-peer encryption is very easy to implement.

You can just use `UpdateCodec` to update serialiser and cryptograph.
Be careful! Both `ws.Service` and `ws.Connection` has this method. `Service` one allows you to specify the default communication codec, and another one allows you to update the communication codec.

Try to implement the proprietary PKI-based ECDH encryption strategy by leveraging the features of WebSocket.