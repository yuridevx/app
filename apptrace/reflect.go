package apptrace

import (
	"runtime"
	"strings"
)

var FunctionToName func(p uintptr) string

func init() {
	FunctionToName = func(p uintptr) string {
		fn := runtime.FuncForPC(p)
		if fn == nil {
			return "unknown function"
		}
		name := fn.Name()
		sp := strings.Split(name, "/")
		nameRaw := sp[len(sp)-1]
		namesp := strings.Split(nameRaw, ".")
		name = namesp[len(namesp)-1]
		finalName := strings.Replace(name, "-fm", "", -1)
		return finalName
	}
}
