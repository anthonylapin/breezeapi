package server

import (
	"fmt"
	"regexp"
	"strings"
)

type RequestHandler = func(ctx Context) HttpResponse

type Router struct {
	handlersPerMethod map[string]map[string]RequestHandler
}

type RoutersRegistry struct {
	Routers []Router
}

func NewRouter() Router {
	return Router{
		handlersPerMethod: make(map[string]map[string]RequestHandler),
	}
}

func (router *Router) addHandler(method string, path string, handler RequestHandler) {
	if router.handlersPerMethod == nil {
		router.handlersPerMethod = make(map[string]map[string]RequestHandler)
	}

	if router.handlersPerMethod[method] == nil {
		router.handlersPerMethod[method] = make(map[string]RequestHandler)
	}

	router.handlersPerMethod[method][path] = handler
	fmt.Println("Registered new handler for ", method, path)
}

func (router *Router) Get(path string, handler RequestHandler) {
	router.addHandler("GET", path, handler)
}

func (router *Router) Post(path string, handler RequestHandler) {
	router.addHandler("POST", path, handler)
}

func (router *Router) Put(path string, handler RequestHandler) {
	router.addHandler("PUT", path, handler)
}

func (router *Router) Delete(path string, handler RequestHandler) {
	router.addHandler("DELETE", path, handler)
}

func extractPathParameter(input string) (bool, string) {
	re := regexp.MustCompile(`^\{([^{}]+)\}$`)

	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		return true, matches[1]
	}

	return false, ""
}

func (router Router) findHandler(request HttpRequest) *RequestHandler {
	handlersForMethod, handlersExistForMethod := router.handlersPerMethod[request.Method]

	if (!handlersExistForMethod) {
		fmt.Println("Not found handler for method", request.Method)
		return nil
	}

	requestPathParts := strings.Split(request.Path, "/")

	for handlerPath, handler := range handlersForMethod {
		pathParts := strings.Split(handlerPath, "/")

		if len(requestPathParts) != len(pathParts) {
			continue
		}

		samePath := true

		for i := 0; i < len(requestPathParts); i++ {
			requestPathPart := requestPathParts[i]
			pathPart := pathParts[i]

			isPathParam, pathParam := extractPathParameter(pathPart)

			if isPathParam {
				request.PathParams[pathParam] = requestPathPart
			} else if requestPathPart != pathPart {
				samePath = false
				break
			}
		}

		if samePath {
			return &handler
		}
	}

	return nil
}

func (registry RoutersRegistry) findHandler(request HttpRequest) *RequestHandler {
	for _, router := range registry.Routers {
		maybeRequestHandler := router.findHandler(request)

		if maybeRequestHandler != nil {
			return maybeRequestHandler
		}
	}

	fmt.Println("Not found a router")
	return nil
}

func (registry *RoutersRegistry) addRouter(router Router) {
	registry.Routers = append(registry.Routers, router)
}