package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/cmd"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func listener(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in handleConnection", r)
			cmd.ServerSideError(conn)
		}
		conn.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			cmd.UnknownError(conn)
			return
		}
		fmt.Printf("Received %d bytes: %s\n", n, buf[:n])

		val, err := resp.NewReader(strings.NewReader(string(buf[:n]))).Read()
		if err != nil {
			fmt.Println("Error parsing:", err)
			cmd.UnknownError(conn)
			return
		}

		_type := strings.ToLower(val.Array[0].Bulk)
		switch _type {
		case cmd.PING_NAME:
			{
				cmd.Ping(conn)
			}
		case cmd.ECHO_NAME:
			{
				cmd.Echo(conn, val)
			}
		case cmd.SET_NAME:
			{
				cmd.Set(conn, val)
			}
		case cmd.GET_NAME:
			{
				cmd.Get(conn, val)
			}
		default:
			{
				cmd.UnknownCommandError(conn)
			}
		}

	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	listener(l)
}
