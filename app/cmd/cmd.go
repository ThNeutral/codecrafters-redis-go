package cmd

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

const (
	PING_NAME = "ping"
	ECHO_NAME = "echo"
	SET_NAME  = "set"
	GET_NAME  = "get"
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
	key := val.Array[1].Bulk
	value := val.Array[2].Bulk
	st.Set(key, value)
	_, err := conn.Write([]byte(rw.Defaults.OK()))
	return err
}

func Get(conn net.Conn, val resp.Value) error {
	key := val.Array[1].Bulk
	str := st.Get(key)
	if str == "" {
		_, err := conn.Write([]byte(rw.Defaults.NULL()))
		return err
	}
	value := resp.Value{
		Type: resp.BULK_NAME,
		Bulk: str,
	}
	respString := rw.Encode(value)
	_, err := conn.Write([]byte(respString))
	return err
}
