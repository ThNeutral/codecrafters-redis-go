package resp

import "fmt"

type RespWriter struct{}

func NewWriter() *RespWriter {
	return &RespWriter{}
}

func (rw *RespWriter) Encode(value Value) string {
	switch value.Type {
	case BULK_NAME:
		{
			return rw.encodeBulk(value)
		}
	default:
		{
			fmt.Printf("Unknown type: %v", string(value.Type))
			return ""
		}
	}
}

func (rw *RespWriter) encodeBulk(value Value) string {
	str := "$"
	str += fmt.Sprintf("%v", len(value.Bulk))
	str += CRLF
	str += value.Bulk
	str += CRLF
	return str
}
