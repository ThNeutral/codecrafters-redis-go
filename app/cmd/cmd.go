package cmd

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	PING_NAME = "PING"
	ECHO_NAME = "ECHO"
)

func Ping(conn net.Conn) error {
	_, err := conn.Write([]byte("+PONG\r\n"))
	return err
}

func Echo(conn net.Conn, val resp.Value) error {
	rw := resp.NewWriter()
	respString := rw.Encode(val)
	_, err := conn.Write([]byte(respString))
	return err
}
