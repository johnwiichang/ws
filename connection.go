package ws

import (
	"errors"
	"fmt"

	"golang.org/x/net/websocket"
)

type (
	//Connection Customized websocket connection definition based on websocket.Conn which provided by official Google golang tool package
	Connection struct {
		Identifier string
		conn       *websocket.Conn
		closed     bool
		codec      *websocket.Codec
		onclosing  []func()
	}
)

//Send 发送
func (conn *Connection) Send(obj interface{}) (err error) {
	if conn.closed {
		err = errors.New("connection has beed closed for a while")
	} else {
		err = conn.codec.Send(conn.conn, obj)
		if err != nil {
			err = errors.New(`an error has been occurred while sending message coz: ` + err.Error())
			if erroronclose := conn.Close(); erroronclose != nil {
				err = errors.New(err.Error() + " and close operation reported an extra error: " + erroronclose.Error())
			}
		}
	}
	return
}

//Receive 获取请求
func (conn *Connection) Receive(obj interface{}) (err error) {
	if conn.closed {
		err = errors.New("connection has beed closed for a while")
	} else {
		err = conn.codec.Receive(conn.conn, obj)
		if err != nil {
			err = errors.New(`an error has been occurred while receiving message coz: ` + err.Error())
			if erroronclose := conn.Close(); erroronclose != nil {
				err = errors.New(err.Error() + " and close operation reported an extra error: " + erroronclose.Error())
			}
		}
	}
	return
}

//RegisterClosingFunc 注册关闭函数
func (conn *Connection) RegisterClosingFunc(handler func(), identifier ...string) {
	conn.onclosing = append(conn.onclosing, handler)
}

//Close 关闭连接
func (conn *Connection) Close() (err error) {
	if !conn.closed {
		err = conn.conn.Close()
		conn.closed = true
		func() {
			defer func() {
				if ex := recover(); ex != nil {
					err = fmt.Errorf("%s", ex)
				}
			}()
			for _, handler := range conn.onclosing {
				handler()
			}
		}()
	}
	return
}

//UpdateCodec 更新连接中使用的译码器
func (conn *Connection) UpdateCodec(codec *websocket.Codec) {
	conn.codec = codec
}
