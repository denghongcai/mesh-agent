package protocol

type Request interface {
	GetData() interface{}
}