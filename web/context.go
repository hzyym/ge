package web

import (
	"encoding/json"
	"gz/web/binding"
	"mime/multipart"
	"net/http"
)

type Context struct {
	Request  *http.Request
	Write    http.ResponseWriter
	index    int
	handle   Handlers
	param    map[string]string
	cxtParam map[string]interface{}

	status int
	engine *Engine
}

func newContext(req *http.Request, write http.ResponseWriter, query map[string]string) *Context {
	return &Context{
		Request:  req,
		Write:    write,
		param:    query,
		cxtParam: make(map[string]interface{}),
	}
}
func (c *Context) write(code int, data []byte) {
	c.status = code
	c.Write.WriteHeader(code)
	_, _ = c.Write.Write(data)
}
func (c *Context) Status(code int) {
	c.status = code
	c.Write.WriteHeader(code)
}
func (c *Context) String(code int, str string) {
	c.write(code, []byte(str))
}
func (c *Context) Query(key string) (string, bool) {
	value, ok := c.param[key]
	return value, ok
}
func (c *Context) QueryForm(key string) string {
	return c.Request.FormValue(key)
}
func (c *Context) Set(key string, data interface{}) {
	c.cxtParam[key] = data
}
func (c *Context) Get(key string) (interface{}, bool) {
	p, ok := c.cxtParam[key]
	return p, ok
}
func (c *Context) Next() {
	c.index++
	l := len(c.handle)

	for ; c.index < l; c.index++ {
		c.handle[c.index](c)
	}
}
func (c *Context) Abort() {
	c.index = 90
}
func (c *Context) GetWriteStatus() int {
	if c.status <= 0 {
		return http.StatusOK
	}
	return c.status
}
func (c *Context) Html(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplate.ExecuteTemplate(c.Write, name, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
func (c *Context) SetHeader(key, value string) {
	c.Write.Header().Set(key, value)
}
func (c *Context) ico() {
	c.File("./static/favicon.ico")
}
func (c *Context) File(path string) {
	http.ServeFile(c.Write, c.Request, path)
}
func (c *Context) Json(code int, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		c.write(http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	c.SetHeader("Content-Type", "application/json")
	c.write(code, buf)
}
func (c *Context) FormParam(key string) string {
	return c.Request.FormValue(key)
}
func (c *Context) PostParam(key string) string {
	return c.Request.PostFormValue(key)
}
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(int64(c.engine.fileSiz)); err != nil {
			return nil, err
		}
	}
	f, h, e := c.Request.FormFile(name)
	_ = f.Close()
	return h, e
}

func (c *Context) RequestHeader(key string) string {
	return c.Request.Header.Get(key)
}
func (c *Context) RequestType() string {
	return c.RequestHeader("Content-Type")
}
func (c *Context) Binding(obj interface{}) error {
	bind := binding.NewBind(c.Request.Method, c.RequestType())
	return bind.Bind(c.Request, obj)
}
func (c *Context) BindingJSON(obj interface{}) error {
	return binding.NewBindType(binding.JSONBind).Bind(c.Request, obj)
}
func (c *Context) BindingXML(obj interface{}) error {
	return binding.NewBindType(binding.XMLBind).Bind(c.Request, obj)
}
func (c *Context) PostForm(key string) string {
	return c.Request.PostFormValue(key)
}
