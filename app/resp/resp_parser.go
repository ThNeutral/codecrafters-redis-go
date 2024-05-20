package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	CRLF = "\r\n"

	ARRAY   = '*'
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'

	ARRAY_NAME   = "array"
	STRING_NAME  = "string"
	ERROR_NAME   = "error"
	INTEGER_NAME = "integer"
	BULK_NAME    = "bulk"
)

type Value struct {
	Type  string
	Str   string
	Num   int
	Error error
	Bulk  string
	Array []Value
}

type RespReader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *RespReader {
	return &RespReader{reader: bufio.NewReader(rd)}
}

func (rr *RespReader) readLine() (line []byte, n int, err error) {
	for {
		b, err := rr.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (rr *RespReader) readInteger() (x int, n int, err error) {
	line, n, err := rr.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (rr *RespReader) Read() (Value, error) {
	_type, err := rr.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return rr.readArray()
	case BULK:
		return rr.readBulk()
	case STRING:
		return rr.readString()
	case INTEGER:
		return rr.readIntegerType()
	case ERROR:
		return rr.readError()
	default:
		fmt.Printf("Unknown type during parsing: %v\n", string(_type))
		return Value{}, nil
	}
}

func (rr *RespReader) readArray() (Value, error) {
	v := Value{}
	v.Type = ARRAY_NAME

	len, _, err := rr.readInteger()
	if err != nil {
		return v, err
	}

	v.Array = make([]Value, 0)
	for range len {
		val, err := rr.Read()
		if err != nil {
			return v, err
		}

		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (rr *RespReader) readBulk() (Value, error) {
	v := Value{}
	v.Type = BULK_NAME

	len, _, err := rr.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	rr.reader.Read(bulk)

	v.Bulk = string(bulk)

	rr.readLine()

	return v, nil
}

func (rr *RespReader) readString() (Value, error) {
	v := Value{}
	v.Type = STRING_NAME

	line, _, err := rr.readLine()
	if err != nil {
		return v, err
	}

	v.Str = string(line)

	return v, nil
}

func (rr *RespReader) readIntegerType() (Value, error) {
	v := Value{}
	v.Type = INTEGER_NAME

	line, _, err := rr.readLine()
	if err != nil {
		return v, err
	}

	i, err := strconv.Atoi(string(line))
	if err != nil {
		return v, err
	}

	v.Num = i

	return v, nil
}

func (rr *RespReader) readError() (Value, error) {
	v := Value{}
	v.Type = ERROR_NAME

	line, _, err := rr.readLine()
	if err != nil {
		return v, err
	}

	v.Error = errors.New(string(line))

	return v, nil
}
