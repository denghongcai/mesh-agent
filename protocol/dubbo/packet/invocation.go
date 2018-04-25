package packet

type Invocation struct {
	methodName     []byte
	args           interface{}
	attachments    map[string]interface{}
	argTypesString []byte
}

func NewInvocation(methodName []byte, args interface{}, attachments map[string]interface{}) *Invocation {
	if attachments[DUBBO_VERSION_KEY] == "" {
		attachments[DUBBO_VERSION_KEY] = DUBBO_VERSION
	}
	return &Invocation{
		methodName:  methodName,
		args:        args,
		attachments: attachments,
	}
}

func (i *Invocation) GetAttachments() map[string]interface{} {
	return i.attachments
}

func (i *Invocation) GetMethodName() []byte {
	return i.methodName
}

func (i *Invocation) SetArgTypesString(s []byte) {
	i.argTypesString = s
}

func (i *Invocation) GetArgTypesString() []byte {
	return i.argTypesString
}

func (i *Invocation) GetArgs() interface{} {
	return i.args
}
