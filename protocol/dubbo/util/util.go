package util

import (
	"reflect"
	"bytes"
)

func getJavaClassDesc(val interface{}) string {
	if val == nil {
		return "Ljava/lang/Object;"
	}

	if reflect.TypeOf(val).Kind() == reflect.String {
		return "Ljava/lang/String;"
	}
	// TODO
	return ""
}

func GetJavaArgsDesc(args interface{}) string {
	var buf bytes.Buffer
	argsArr := args.([]interface{})
	for _, arg := range argsArr {
		buf.WriteString(getJavaClassDesc(arg))
	}
	return buf.String()
}
