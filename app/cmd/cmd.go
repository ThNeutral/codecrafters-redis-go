package cmd

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

const (
	PING_NAME = "ping"
	ECHO_NAME = "echo"
	SET_NAME  = "set"
	GET_NAME  = "get"

	EXPIRY_NAME = "px"
)

var (
	st = storage.NewStorage()
	rw = resp.NewWriter()
)

func Ping(conn net.Conn) error {
	_, err := conn.Write([]byte("+PONG\r\n"))
	return err
}

func Echo(conn net.Conn, val resp.Value) error {
	respString := rw.Encode(val.Array[1])
	_, err := conn.Write([]byte(respString))
	return err
}

func Set(conn net.Conn, val resp.Value) error {
	key := val.Array[1]
	value := val.Array[2]

	if len(val.Array) == 5 && strings.ToLower(val.Array[3].Bulk) == EXPIRY_NAME {
		_int, err := strconv.Atoi(val.Array[4].Bulk)
		if err != nil {
			return err
		}
		st.SetWithExpiry(key, value, time.Duration(_int)*time.Millisecond)
	} else {
		st.Set(key, value)
	}

	_, err := conn.Write([]byte(rw.Defaults.OK()))
	return err
}

func Get(conn net.Conn, val resp.Value) error {
	key := val.Array[1]
	storeValue := st.Get(key)
	if storeValue == nil {
		_, err := conn.Write([]byte(rw.Defaults.NULL()))
		return err
	}
	respString := rw.Encode(*storeValue)
	_, err := conn.Write([]byte(respString))
	return err
}
