package web

import (
	"fmt"
)

type Debug struct {
	isDebug bool
}

func (d *Debug) DebugPrint(event string, content ...interface{}) {
	d.debugOutput("[%s] %s\n", event, content)
}
func (d *Debug) debugOutput(format string, content ...interface{}) {
	if d.isDebug {
		fmt.Printf("[GZ-Debug] "+format, content...)
	}
}
