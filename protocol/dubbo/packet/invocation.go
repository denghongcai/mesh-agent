package packet

type Invocation struct {
	methodName string
	args interface{}
	attachments map[string]string
	argTypesString string
}

func NewInvocation(methodName string, args interface{}, attachments map[string]string) *Invocation {
	if attachments[DUBBO_VERSION_KEY] == "" {
		attachments[DUBBO_VERSION_KEY] = DUBBO_VERSION
	}
	return &Invocation{
		methodName:methodName,
		args:args,
		attachments:attachments,
	}
}

func (i *Invocation) GetAttachments() map[string]string {
	return i.attachments
}

func (i *Invocation) GetMethodName() string {
	return i.methodName
}

func (i *Invocation) SetArgTypesString(s string) {
	i.argTypesString = s
}

func (i *Invocation) GetArgTypesString() string {
	return i.argTypesString
}

func (i *Invocation) GetArgs() interface{} {
	return i.args
}