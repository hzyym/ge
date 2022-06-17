package web

import "net/http"

type groupRoute struct {
	engine   *Engine
	basePath string

	middleware []Handle
}
type IRoute interface {
	Group(basePath string) IRoute
	Use(handle ...Handle)
	Get(path string, handle Handle)
	Post(path string, handle Handle)
}

func newGroup() groupRoute {
	return groupRoute{}
}
func (r *groupRoute) Group(basePath string) IRoute {
	return &groupRoute{
		engine:   r.engine,
		basePath: basePath,
	}
}
func (r *groupRoute) Use(handle ...Handle) {
	r.engine.addMiddleware(r.basePath, handle)
}
func (r *groupRoute) Get(path string, handle Handle) {
	r.handle(http.MethodGet, path, handle)
}
func (r *groupRoute) Post(path string, handle Handle) {
	r.handle(http.MethodPost, path, handle)
}
func (r *groupRoute) handle(method, path string, handle_ Handle) {
	if r.basePath != "" {
		path = r.basePath + path
	}
	r.engine.addRoute(method, path, handle_)
}
