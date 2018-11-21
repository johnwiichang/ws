package ws

import (
	"fmt"
	"net/http"

	"github.com/johnwiichang/ws/codec"
	"golang.org/x/net/websocket"
)

type (
	//utilities 工具类声明
	utilities interface {
		DefaultJSONCodec(...codec.CryptoService) *websocket.Codec
		DefaultByteCodec(...codec.CryptoService) *websocket.Codec
		CustomizedCodec(codec.MarshalService, ...codec.CryptoService) *websocket.Codec

		originalCheck(*websocket.Config, *http.Request) error
	}

	util struct{}
)

//Utils 工具类对象
var Utils utilities = util{}

//DefaultCodec 获取默认编码解码器
func (u util) DefaultJSONCodec(crypto ...codec.CryptoService) *websocket.Codec {
	return u.CustomizedCodec(codec.JSONMarshalService, crypto...)
}

func (u util) DefaultByteCodec(crypto ...codec.CryptoService) *websocket.Codec {
	return u.CustomizedCodec(codec.RawMarshalService, crypto...)
}

func (u util) CustomizedCodec(marshalService codec.MarshalService, cryptoService ...codec.CryptoService) *websocket.Codec {
	marshal := marshalService
	if marshal == nil {
		marshal = codec.RawMarshalService
	}
	var crypto codec.CryptoService
	if len(cryptoService) > 0 {
		crypto = cryptoService[0]
	}
	return &websocket.Codec{
		Marshal: func(v interface{}) (data []byte, payloadType byte, err error) {
			obj, isBytes := v.([]byte)
			if isBytes {
				data = obj
			} else if str, isStr := v.(string); isStr {
				data = []byte(str)
			} else {
				data, err = marshal.Marshal(v)
				payloadType = websocket.BinaryFrame
			}
			if crypto != nil {
				data = crypto.Encrypt(data)
			}
			return
		},
		Unmarshal: func(data []byte, payloadType byte, v interface{}) (err error) {
			var decrypted []byte
			if crypto != nil {
				decrypted, err = crypto.Decrypt(data)
				if err != nil {
					return err
				}
			} else {
				decrypted = data
			}
			return marshal.Unmarshal(decrypted, v)
		},
	}
}

func (u util) originalCheck(config *websocket.Config, req *http.Request) (err error) {
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}
