package web

import "sync"

type contextPool struct {
	pool sync.Pool
}

func newContextPoll() *contextPool {
	p := &contextPool{}
	p.pool.New = func() any {
		return newContext(nil, nil, nil)
	}
	return p
}
func (c *contextPool) Pop() *Context {
	return c.pool.Get().(*Context)
}
func (c *contextPool) Push(cxt *Context) {
	c.pool.Put(cxt)
}
