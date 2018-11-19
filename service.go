package ws

import (
	"errors"
	"net/http"
	"reflect"

	"golang.org/x/net/websocket"
)

var (
	connectionpool = map[string]*Connection{}
)

type (
	//Service WebSocket 服务对象
	Service struct {
		s                *websocket.Server
		handlers         map[string]Handler
		modelConstructor func() interface{}
		codec            *websocket.Codec
		adapter          Adapter
	}

	//Request 请求实现
	Request interface {
		Action() string
		Body() interface{}
	}

	//Handler 句柄声明
	Handler func(*Connection, interface{})

	//Adapter 适配器声明
	Adapter func(*Connection, []byte)
)

//NewWebSocketService 创建新的 WebSocket 服务
func NewWebSocketService(reqObj Request, nonbrowserClient ...bool) *Service {
	originCheck := Utils.originalCheck
	if len(nonbrowserClient) > 0 && nonbrowserClient[0] {
		originCheck = func(config *websocket.Config, req *http.Request) (err error) { return }
	}
	s := &websocket.Server{
		Handler:   nil,
		Handshake: originCheck,
	}
	service := &Service{
		s:        s,
		handlers: make(map[string]Handler, 0),
	}
	service.codec = Utils.DefaultJSONCodec()
	modelRef := reflect.Indirect(reflect.ValueOf(reqObj)).Type()
	service.modelConstructor = func() interface{} { return reflect.New(modelRef).Interface() }
	service.s.Handler = service.ServeWebSocket
	return service
}

//NewWebSocketAdapter 创建 WebSocket 适配器
func NewWebSocketAdapter(handler Adapter, nonbrowserClient ...bool) *Service {
	originCheck := Utils.originalCheck
	if len(nonbrowserClient) > 0 && nonbrowserClient[0] {
		originCheck = func(config *websocket.Config, req *http.Request) (err error) { return }
	}
	s := &websocket.Server{
		Handler:   nil,
		Handshake: originCheck,
	}
	service := &Service{
		s:       s,
		adapter: handler,
	}
	service.codec = Utils.DefaultByteCodec()
	service.modelConstructor = func() interface{} { return &[]byte{} }
	service.s.Handler = service.ServeWebSocket
	return service
}

//ServeWebSocket 实现 ServeWebSocket 接口
func (s *Service) ServeWebSocket(connection *websocket.Conn) {
	defer connection.Close()
	conn := &Connection{
		conn:  connection,
		codec: s.codec,
	}
	for !conn.closed {
		request := s.modelConstructor()
		if err := conn.Receive(request); err != nil {
			conn.Close()
			break
		}
		if s.adapter == nil {
			handler, existed := s.handlers[request.(Request).Action()]
			if existed {
				handler(conn, request.(Request).Body())
			} else {
				conn.Close()
				break
			}
		} else {
			s.adapter(conn, *request.(*[]byte))
		}
	}
}

//UpdateDefaultCodec 更新 Codec
func (s *Service) UpdateDefaultCodec(codec *websocket.Codec) {
	s.codec = codec
}

//ServeHTTP 实现 ServeHTTP 接口
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.s.ServeHTTP(w, r)
}

//RegisterEndpoint 注册 Endpoint
func (s *Service) RegisterEndpoint(endpoint string, handler Handler) (err error) {
	if s.adapter != nil {
		panic(`the endpoint cannot workwith []byte request structure, try using 'RegisterAdapter' instead`)
	} else {
		if _, existed := s.handlers[endpoint]; existed {
			err = errors.New("repeated handler registration: " + endpoint)
		}
		s.handlers[endpoint] = handler
	}
	return
}
