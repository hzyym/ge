package binding

import "net/http"

var (
	form  = formBinding{}
	jsons = jsonBinding{}
	xmls  = xmlBinding{}
)

const (
	JSONBind BindType = iota
	FORMBind
	XMLBind
)

type BindType int
type Binding interface {
	Bind(req *http.Request, obj interface{}) error
}

func NewBind(method, reqType string) Binding {
	if method == http.MethodGet {
		return form
	}
	switch reqType {
	case "application/json":
		return jsons
	case "application/xml", "text/xml":
		return xmls
	}
	return form
}
func NewBindType(types BindType) Binding {
	switch types {
	case JSONBind:
		return jsons
	case FORMBind:
		return form
	case XMLBind:
		return xmls
	}
	return form
}
