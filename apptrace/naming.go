package apptrace

import (
	"github.com/yuridevx/app/extension"
	"runtime"
	"strings"
)

func GetNameForPc(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
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

type NamingStrategy func(call extension.CallType, part extension.Part) string

var DefaultNamingStrategy = func(call extension.CallType, part extension.Part) string {
	return GetNameForPc(part.GetHandler())
}
