package web

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type route struct {
	node            map[string]*node
	groupMiddleware map[string]Handlers

	notRouteFun Handle

	staticFile http.Handler
}

func newRoute() *route {
	return &route{node: make(map[string]*node), groupMiddleware: make(map[string]Handlers)}
}
func (r *route) parsePattern(pattern string) (part []string) {
	if pattern == "/" {
		part = append(part, "/")
		return part
	}
	val := strings.Split(pattern, "/")
	for _, v := range val {
		if v != "" {
			part = append(part, v)
		}
	}
	return
}
func (r *route) addRoute(methodType string, path string, handle Handle) {
	nodes, ok := r.node[methodType]
	part := r.parsePattern(path)
	if !ok {
		nodes = &node{
			pattern:  "/",
			part:     "/",
			children: nil,
			isWild:   false,
			types:    root,
			handle:   nil,
		}
		r.node[methodType] = nodes
	}
	nodes.insert(path, part, 0, handle)
}
func (r *route) getRoute(methodType string, path string) (Handlers, map[string]string) {
	nodes, ok := r.node[methodType]
	if !ok {
		return Handlers{r.NotRoute404()}, nil
	}
	var handleAll Handlers
	nodes, params := nodes.search(r.parsePattern(path), 0, func(part string) {
		h, ok := r.groupMiddleware[part]
		if !ok {
			return
		}
		handleAll = append(handleAll, h...)
	})
	var h Handle
	if nodes == nil || nodes.handle == nil {
		h = r.NotRoute404()
	} else {
		h = nodes.handle
	}
	handleAll = append(handleAll, h)
	return handleAll, params
}
func (r *route) NotRoute404() Handle {
	if r.notRouteFun == nil {
		return func(c *Context) {
			c.String(http.StatusNotFound, "gz web page 404")
		}
	}
	return r.notRouteFun
}
func (r *route) SetNotRoute(handle Handle) {
	r.notRouteFun = handle
}
func (r *route) addGroupMiddleware(basePath string, handle Handlers) {
	if basePath == "/" {
		r.groupMiddleware[basePath] = append(r.groupMiddleware[basePath], handle...)
		return
	}
	part := strings.Split(basePath, "/")

	r.groupMiddleware[part[len(part)-1]] = handle
}
func (r *route) printRouter() {
	var path string

	fmt.Println("Route ALL:")
	path += "/"
	for reqType, nodes := range r.node {
		if nodes.types != root {
			r.routeInfo(reqType, nodes, "/")
		} else if nodes.types == root && nodes.handle != nil {
			r.println(reqType, path, r.FunName(nodes.handle))
		}
		if nodes.children != nil {
			r.routeInfo(reqType, nodes, "/")
		}
	}
}
func (r *route) routeInfo(reqType string, nodes *node, path string) {
	if nodes.children == nil {
		r.println(reqType, path, r.FunName(nodes.handle))
		return
	} else {
		for _, childNode := range nodes.children {
			if childNode.children != nil {
				r.routeInfo(reqType, childNode, path+childNode.part+"/")
				continue
			} else {
				r.println(reqType, path+childNode.part, r.FunName(nodes.handle))
			}
		}
	}
}
func (r *route) FunName(handle Handle) string {
	return runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name()
}
func (r *route) println(reqType, path, funName string) {
	fmt.Printf("\t[%s] %s     --->%s\n", reqType, path, funName)
}
func (r *route) static(roots, path string) Handle {
	f := http.Dir(roots)
	r.staticFile = http.StripPrefix(path, http.FileServer(f))

	return func(c *Context) {
		name, _ := c.Query("file")
		if _, err := f.Open(name); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		r.staticFile.ServeHTTP(c.Write, c.Request)
	}
}
