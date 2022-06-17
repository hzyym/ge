package binding

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type testForm struct {
	I int    `form:"i"`
	Z string `form:"z"`
}

func TestForm(t *testing.T) {
	req := new(http.Request)

	req.Form = make(url.Values)
	req.Form.Set("i", "20")
	req.Form.Set("z", "he")

	v := make(map[string]string)
	_ = mapForm(req, reflect.ValueOf(&v), "form")

	fmt.Println(v)

}
