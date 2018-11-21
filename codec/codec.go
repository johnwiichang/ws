package codec

type (
	CryptoService interface {
		Encrypt([]byte) []byte
		Decrypt([]byte) ([]byte, error)
	}

	MarshalService interface {
		Marshal(interface{}) ([]byte, error)
		Unmarshal([]byte, interface{}) error
	}
)
