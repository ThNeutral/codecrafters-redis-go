package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

var (
	rw = resp.NewWriter()
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Printf("Error ocured while connecting: %v \n", err.Error())
		os.Exit(1)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("client > ")
		input, _ := reader.ReadString('\n')
		input = string(input[:len(input)-2])

		vArray := resp.Value{}
		vArray.Type = resp.ARRAY_NAME
		splitted := strings.Split(input, " ")
		for _, sp := range splitted {
			vBulk := resp.Value{}
			vBulk.Type = resp.BULK_NAME
			vBulk.Bulk = sp
			vArray.Array = append(vArray.Array, vBulk)
		}

		_, err = conn.Write([]byte(rw.Encode(vArray)))
		if err != nil {
			fmt.Printf("Error writing bytes: %v\n", err.Error())
			continue
		}

		buff := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("Error reading bytes: %v\n", err.Error())
			continue
		}

		fmt.Println("\n", string(buff[:n]))
	}
}
