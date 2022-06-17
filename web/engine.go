package web

import (
	"html/template"
	"net/http"
)

type EngineInterface interface {
	Run(address string) error
	NotRoute(handle Handle)
	Use(handle ...Handle)
	Group(basePath string) IRoute
	Get(path string, handle Handle)
	Post(path string, handle Handle)
	Static(roots, path string)
	LoadHtmlGlob(pattern string)
	Delim(left, right string)
	SetUploadFileSiz(siz int)
}

func New() EngineInterface {
	return newEngine()
}
func Default() EngineInterface {
	engine := newEngine()

	engine.Use(recovery(), requestLog())
	return engine
}
func newEngine() *Engine {
	engine := &Engine{
		newGroup(),
		newRoute(),
		newContextPoll(),
		&Debug{true},
		nil,
		"",
		32 << 20,
	}

	engine.groupRoute.engine = engine
	return engine
}

type Handle func(c *Context)
type Handlers []Handle
type Engine struct {
	groupRoute
	route *route

	pool *contextPool

	debug *Debug

	htmlTemplate *template.Template
	htmlPattern  string

	fileSiz int
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h, params := e.route.getRoute(request.Method, request.URL.Path)
	cxt := e.pool.Pop()
	cxt.index = -1
	cxt.Write = writer
	cxt.param = params
	cxt.Request = request
	cxt.handle = h
	cxt.engine = e

	cxt.Next()

	e.pool.Push(cxt)
}
func (e *Engine) Run(address string) error {
	e.debug.DebugPrint("Run-Server", "Start Address"+address)

	if e.debug.isDebug {
		e.route.printRouter()
	}
	err := http.ListenAndServe(address, e)
	if err != nil {
		e.debug.DebugPrint("Run Server Fail", err.Error())
	}
	return err
}
func (e *Engine) NotRoute(handle Handle) {
	e.route.SetNotRoute(handle)
}
func (e *Engine) addRoute(method, path string, handle Handle) {
	e.route.addRoute(method, path, handle)
}
func (e *Engine) addMiddleware(basePath string, handle Handlers) {
	e.route.addGroupMiddleware(basePath, handle)
}
func (e *Engine) Use(handle ...Handle) {
	e.addMiddleware("/", handle)
}
func (e *Engine) Static(roots, path string) {
	h := e.route.static(roots, path)
	e.Get(path+"/*file", h)
	e.Get("/favicon.ico", func(c *Context) {
		c.ico()
	})
}

func (e *Engine) LoadHtmlGlob(pattern string) {
	e.htmlPattern = pattern
	e.htmlTemplate = template.Must(template.New("").ParseGlob(pattern))
}
func (e *Engine) Delim(left, right string) {
	if e.htmlTemplate == nil {
		return
	}
	e.htmlTemplate = template.Must(template.New("").Delims(left, right).ParseGlob(e.htmlPattern))
}
func (e *Engine) SetUploadFileSiz(siz int) {
	e.fileSiz = siz
}
