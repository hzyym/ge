package binding

import "net/http"

var (
	form = formBinding{}
)

type Binding interface {
	Bind(req *http.Request, obj interface{}) error
}

func NewBing(method, reqType string) Binding {
	if method == http.MethodGet {
		return form
	}
	switch reqType {
	case "application/json":

	}
	return form
}
