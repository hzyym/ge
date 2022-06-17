package web

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func recovery() Handle {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				trace(fmt.Sprintf("%s", err))

				if c.GetWriteStatus() <= 0 {
					c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
			}
		}()

		c.Next()
	}
}
func trace(message string) {
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)
	var strBuf strings.Builder
	strBuf.WriteString(message + " -> err")
	for _, val := range pc[:n] {
		fun := runtime.FuncForPC(val)

		file, lien := fun.FileLine(val)

		strBuf.WriteString(fmt.Sprintf("\n\t%s:%d funName:%s", file, lien, fun.Name()))
	}
	fmt.Println(strBuf.String())
}
