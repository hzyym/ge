package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type Json struct {
	I int    `json:"i"`
	Z string `json:"z"`
}

func TestJson(t *testing.T) {
	req := new(http.Request)
	var testData Json

	testData.I = 20
	testData.Z = "he"
	buf, _ := json.Marshal(testData)
	req.Body = ioutil.NopCloser(bytes.NewReader(buf))

	var tmp Json
	err := unmarshal(req, reflect.ValueOf(tmp))

	fmt.Println(err)
	fmt.Println(tmp)
}
