package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var STATUS_CODE_TO_MESSAGE_MAP = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	200: "OK",
	201: "Created",
	204: "No Content",
	301: "Moved Permanently",
	302: "Found",
	304: "Not Modified",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

type HttpResponse struct {
	HttpVersion string
	Code        int
	Headers     map[string]string
	Data        []byte
}

func (response *HttpResponse) setData(data []byte, contentType string) {
	if contentEncoding, contentEncodingExists := response.Headers["Content-Encoding"]; contentEncodingExists && compressionSupported(contentEncoding) {
		compressedData, _ := getCompressedData(contentEncoding, data)
		response.Data = compressedData
	} else {
		response.Data = data
	}

	response.Headers["Content-Length"] = strconv.Itoa(len(response.Data))
	response.Headers["Content-Type"] = contentType
}

func (response *HttpResponse) setEncoding(ctx Context) {
	clientAcceptedEncodings, hasClientAcceptedEncodings := ctx.Request.Headers["Accept-Encoding"]

	if !hasClientAcceptedEncodings {
		return
	}

	for _, acceptedEncoding := range strings.Split(clientAcceptedEncodings, ",") {
		if trimmedEncoding := strings.TrimSpace(acceptedEncoding); compressionSupported(trimmedEncoding) {
			response.Headers["Content-Encoding"] = trimmedEncoding
		}
	}
}

func NewResponse(ctx Context, statusCode int) HttpResponse {
	response := HttpResponse{
		HttpVersion: ctx.Request.HttpVersion,
		Code:        statusCode,
		Headers:     make(map[string]string),
	}

	response.setEncoding(ctx)

	return response
}

func NotFoundResponse(ctx Context) HttpResponse {
	return NewResponse(ctx, 404)
}

func OkResponse(ctx Context) HttpResponse {
	return NewResponse(ctx, 200)
}

func OkResponseWithText(ctx Context, text string) HttpResponse {
	response := OkResponse(ctx)
	response.setData([]byte(text), "text/plain")
	return response
}

func OkResponseWithJson(ctx Context, data any) HttpResponse {
	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Failed to encode as JSON", err)
		return InternalServerErrorResponse(ctx)
	}

	response := OkResponse(ctx)
	response.setData([]byte(jsonData), "application/json")
	return response
}

func OkResponseWithFile(ctx Context, file []byte) HttpResponse {
	response := OkResponse(ctx)
	response.setData(file, "application/octet-stream")
	return response
}

func BadRequestResponse(ctx Context) HttpResponse {
	return NewResponse(ctx, 400)
}

func InternalServerErrorResponse(ctx Context) HttpResponse {
	return NewResponse(ctx, 500)
}

func CreatedResponse(ctx Context) HttpResponse {
	return NewResponse(ctx, 201)
}

