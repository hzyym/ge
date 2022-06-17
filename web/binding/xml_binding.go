package binding

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"reflect"
)

type xmlBinding struct {
}

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	return xmlUnmarshal(req, reflect.ValueOf(obj))
}
func xmlUnmarshal(req *http.Request, obj reflect.Value) error {
	if req.Method == http.MethodGet {
		return paresForm(req, obj, "xml")
	}
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}
	return xml.Unmarshal(buf, obj.Interface())
}
