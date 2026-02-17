package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Starting Redis Server...")

	incChan := make(chan struct{})
	decChan := make(chan struct{})

	go counterManager(incChan, decChan)

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

		incChan <- struct{}{}
		go acceptConnection(conn, decChan)
	}
}

func acceptConnection(conn net.Conn, decChan chan struct{}) {
	fmt.Println("Connection Established...", conn.RemoteAddr())
	defer func() {
		conn.Close()
		decChan <- struct{}{}
	}()

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

		resp, err := parseInput(b[:n])
		if err != nil {
			fmt.Println("error parsing input")
			return
		}

		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println("Error Sending Response: ", err.Error())
			return
		}
	}
}

func counterManager(incChan, decChan chan struct{}) {
	counter := 0

	for {
		select {
		case <-incChan:
			counter++
			fmt.Printf("Connection opened. Active connections: %d\n", counter)
		case <-decChan:
			counter--
			fmt.Printf("Connection closed. Active connections: %d\n", counter)
		}
	}
}

func parseInput(input []byte) ([]byte, error) {

	if len(input) == 0 {
		return nil, errors.New("empty input")
	}

	const CarriageReturn = '\r'
	const LineFeed = '\n'

	inputAsString := string(input)
	inputType := inputAsString[0]

	if len(input) < 2 {
		return nil, errors.New("malformed input, missing size")
	}

	inputSize, err := strconv.Atoi(string(inputAsString[1]))

	if err != nil {
		return nil, errors.New("error converting input size")
	}

	fmt.Printf("Type: %c, Size: %d\n", inputType, inputSize)
	sizeRemoved := inputAsString[2:]

	cleanInput := strings.TrimSuffix(sizeRemoved, "\r\n")
	cleanInput = strings.TrimPrefix(cleanInput, "\r\n")
	result := strings.Split(cleanInput, "\r\n")

	fmt.Println(result)

	command := result[1]

	// Handle casing
	// for now we decide that we accept mixed-casing and wrong casing of the commands
	command = strings.ToUpper(command)

	switch command {
	case "ECHO":
		value, ok := getAt(result, 3)
		if !ok {
			fmt.Printf("error reading ECHO input, not enough parameters")
			return []byte("+ERROR\r\n"), errors.New("error not enough parameters")
		}
		return buildBulkString(value), nil
	case "PING":
		return []byte("+PONG\r\n"), nil
	default:
		return nil, fmt.Errorf("unrecognized command: %s", command)
	}
}

func getAt(input []string, idx int) (string, bool) {
	if idx < len(input) && idx > 0 {
		return input[idx], true
	}

	return "", false

}

// Bulk String
// $<length>\r\n<data>\r\n
func buildBulkString(text string) []byte {
	return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(text), text)
}
