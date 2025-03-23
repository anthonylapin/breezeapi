package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type HttpRequest struct {
	Method      string
	Path        string
	HttpVersion string
	Headers     map[string]string
	PathParams  map[string]string
	Body        []byte
}

func (request *HttpRequest) parseRequestLine(requestLine string) error {
	args := strings.Split(requestLine, " ")

	if len(args) != 3 {
		return fmt.Errorf("Invalid number of args. Expected: %d, received: %d", 3, len(args))
	}

	request.Method = strings.TrimSpace(args[0])
	request.Path = strings.TrimSpace(args[1])
	request.HttpVersion = strings.TrimSpace(args[2])

	return nil
}

func (request *HttpRequest) addHeader(headerLine string) error {
	header := strings.Split(headerLine, ": ")

	if len(header) != 2 {
		return fmt.Errorf("Failed to parse header line: %s", headerLine)
	}

	request.Headers[strings.TrimSpace(header[0])] = strings.TrimSpace(header[1])

	return nil
}

func NewRequest(connection net.Conn) (HttpRequest, error) {
	request := HttpRequest{Headers: make(map[string]string), PathParams: make(map[string]string)}

	reader := bufio.NewReader(connection)

	requestLine, err := reader.ReadString('\n')

	if err != nil {
		return request, err
	}

	fmt.Println("Received request:", requestLine)
	err = request.parseRequestLine(requestLine)

	if err != nil {
		return request, err
	}

	// Read headers
	for {
		line, err := reader.ReadString('\n')

		if err != nil || line == "\r\n" {
			break
		}

		request.addHeader(line)
	}

	// read body if content length is presented
	contentLengthStr, contentLengthExists := request.Headers["Content-Length"]
	if contentLengthExists {
		contentLength := 0
		fmt.Sscanf(contentLengthStr, "%d", &contentLength)

		request.Body = make([]byte, contentLength)
		_, err := reader.Read(request.Body)

		if err != nil {
			fmt.Println("Error reading body", err)
			return request, err
		}
	}

	return request, nil
}
