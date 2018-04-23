package packet

type Invocation struct {
	methodName string
	iface string
	args interface{}
	argTypesString string
}

func NewInvocation(methodName string, iface string, args interface{}) *Invocation {
	return &Invocation{
		methodName:methodName,
		args:args,
		iface:iface,
	}
}

func (i *Invocation) SetArgTypesString(s string) {
	i.argTypesString = s
}

func (i *Invocation) GetArgTypesString() string {
	return i.argTypesString
}

func (i *Invocation) GetInterface () string {
	return i.iface
}

func (i *Invocation) GetMethodName() string {
	return i.methodName
}

func (i *Invocation) GetArgs() interface{} {
	return i.args
}
