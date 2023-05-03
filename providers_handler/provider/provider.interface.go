package provider

type IProvider interface {
	Call(chainType string, body []byte) ([]byte, error)
}