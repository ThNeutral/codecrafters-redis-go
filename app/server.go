package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) error {
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		return err
	}
	defer conn.Close()

	pongs := strings.Split(string(buff), "\n")

	fmt.Println(buff)

	for _, _ = range pongs {
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn)
	}
}
