package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var (
	_ = net.Listen
	_ = os.Exit
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Server Started: ", l.Addr())
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		go acceptConnection(conn)

	}
}

func acceptConnection(conn net.Conn) {
	fmt.Println("Connection Established...", conn.RemoteAddr())

	for {
		b := make([]byte, 128)
		n, err := conn.Read(b)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client: ", conn.RemoteAddr())
			} else {
				fmt.Println("Error Reading Connection: ", err.Error())
			}

			return
		}
		parseInput(b[:n])

		resp := []byte("+PONG\r\n")
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Error Sending Response: ", err.Error())
			return
		}
	}
}

func parseInput(input []byte) []string {
	if len(input) == 0 {
		return make([]string, 0)
	}

	dataType := input[1]

	switch dataType {
	case '$':
		fmt.Println("Bulk String")
	}

	return []string{}
}
