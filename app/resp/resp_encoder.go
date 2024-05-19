package resp

import (
	"fmt"
	"strconv"
)

type RespWriter struct {
	Defaults *defaultsType
}

func NewWriter() *RespWriter {
	return &RespWriter{
		Defaults: &defaultsType{},
	}
}

func (rw *RespWriter) Encode(value Value) string {
	switch value.Type {
	case BULK_NAME:
		{
			return rw.encodeBulk(value)
		}
	case STRING_NAME:
		{
			return rw.encodeString(value)
		}
	case INTEGER_NAME:
		{
			return rw.encodeInt(value)
		}
	case ARRAY_NAME:
		{
			return rw.encodeArray(value)
		}
	default:
		{
			fmt.Printf("Unknown type: %v\n", string(value.Type))
			return ""
		}
	}
}

func (rw *RespWriter) encodeBulk(value Value) string {
	str := string(BULK)
	str += fmt.Sprintf("%v", len(value.Bulk))
	str += CRLF
	str += value.Bulk
	str += CRLF
	return str
}

func (rw *RespWriter) encodeString(value Value) string {
	str := string(STRING)
	str += value.Str
	str += CRLF
	return str
}

func (rw *RespWriter) encodeInt(value Value) string {
	str := string(INTEGER)
	str += strconv.Itoa(value.Num)
	str += CRLF
	return str
}

func (rw *RespWriter) encodeArray(value Value) string {
	str := string(ARRAY)
	l := len(value.Array)
	str += strconv.Itoa(l)
	str += CRLF
	for _, v := range value.Array {
		str += rw.Encode(v)
		str += CRLF
	}
	return str
}

type defaultsType struct {
}

func (defaults *defaultsType) OK() string {
	rw := NewWriter()
	v := Value{
		Type: STRING_NAME,
		Str:  "OK",
	}
	return rw.Encode(v)
}

func (defaults *defaultsType) NULL() string {
	return "$-1\r\n"
}
