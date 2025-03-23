package server

type Context struct {
	Request HttpRequest
}

func newContext() Context {
	return Context{}
}
