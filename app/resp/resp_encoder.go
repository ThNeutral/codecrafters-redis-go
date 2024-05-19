package resp

import "fmt"

type defaultsType struct {
}

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
	default:
		{
			fmt.Printf("Unknown type: %v", string(value.Type))
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
