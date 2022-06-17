package binding

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

type jsonBinding struct {
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	return unmarshal(req, reflect.ValueOf(obj))
}
func unmarshal(req *http.Request, obj reflect.Value) error {
	if req.Method == http.MethodGet {
		return paresForm(req, obj, "json")
	}
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, obj.Interface())
	return err
}
