package server

import (
	"fmt"
	"net"
)

func sendResponse(connection net.Conn, response HttpResponse) error {
	statusCodeMessage, statusCodeMessageExists := STATUS_CODE_TO_MESSAGE_MAP[response.Code]

	if !statusCodeMessageExists {
		return fmt.Errorf("Failed to find status code message for the %d code", response.Code)
	}

	connection.Write([]byte(fmt.Sprintf("%s %d %s\r\n", response.HttpVersion, response.Code, statusCodeMessage)))
	
	if response.Headers != nil {
		for headerKey, headerValue := range response.Headers {
			connection.Write([]byte(fmt.Sprintf("%s: %s\r\n", headerKey, headerValue)))
		}
	}
	connection.Write([]byte("\r\n")) // End of headers

	connection.Write([]byte(response.Data))

	fmt.Println("Sent response", response)
	return nil
}

func getResponse(ctx Context, routersRegistry RoutersRegistry) HttpResponse {
	requestHandler := routersRegistry.findHandler(ctx.Request)

	if requestHandler == nil {
		return NotFoundResponse(ctx)
	}

	response := (*requestHandler)(ctx)
	return response
}

func handleConnection(ctx Context, connection net.Conn, routersRegistry *RoutersRegistry) {
	fmt.Println("Handling connection")
	defer connection.Close()

	request, err := NewRequest(connection)

	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	ctx.Request = request

	response := getResponse(ctx, *routersRegistry)
	sendResponse(connection, response)
}