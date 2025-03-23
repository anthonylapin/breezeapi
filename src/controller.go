package main

import (
	"breezeapi/src/server"
	"fmt"
	"os"
	"path"
)

type RequestController struct {
	CLIArgs map[string]string
}

func (controller RequestController) ping(ctx server.Context) server.HttpResponse {
	return server.OkResponse(ctx)
}

func (controller RequestController) echo(ctx server.Context) server.HttpResponse {
	return server.OkResponseWithText(ctx, ctx.Request.PathParams["str"])
}

func (controller RequestController) userAgent(ctx server.Context) server.HttpResponse {
	return server.OkResponseWithText(ctx, ctx.Request.Headers["User-Agent"])
}

func (controller RequestController) getFile(ctx server.Context) server.HttpResponse {
	fileDir, fileDirExists := controller.CLIArgs["fileDirectory"]

	if !fileDirExists {
		return server.NotFoundResponse(ctx)
	}

	fileName := ctx.Request.PathParams["fileName"]
	filePath := path.Join(fileDir, fileName)

	fmt.Println(filePath)

	fileData, err := os.ReadFile(filePath)

	if err != nil {
		return server.NotFoundResponse(ctx)
	}

	return server.OkResponseWithFile(ctx, fileData)
}

func (controller RequestController) writeFile(ctx server.Context) server.HttpResponse {
	fileDir, fileDirExists := controller.CLIArgs["fileDirectory"]

	if !fileDirExists {
		return server.NotFoundResponse(ctx)
	}

	fileName := ctx.Request.PathParams["fileName"]
	filePath := path.Join(fileDir, fileName)

	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return server.BadRequestResponse(ctx)
	}

	defer file.Close()

	_, err = file.Write(ctx.Request.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return server.InternalServerErrorResponse(ctx)
	}

	return server.CreatedResponse(ctx)
}
