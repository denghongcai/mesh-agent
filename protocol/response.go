package protocol

type Response interface {
	GetData() interface{}
}