package web

import "fmt"

func requestLog() Handle {
	return func(c *Context) {

		c.Next() //后中间件
		fmt.Println("["+c.Request.Method+"] ", c.GetWriteStatus(), " ---> "+c.Request.URL.Path)
	}
}
