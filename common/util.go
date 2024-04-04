package common

import (
	"runtime"
	"strings"
)

// Trace 는 현재 호출된 함수 라인등 정보를 가져옴
func Trace() (file string, line int, functionName string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "?"
	}

	return file, line, fn.Name()
}

// FunctionName 현재 호출된 경로를 제외한 패키지 포함 function 이름 조회
// ex) common.FunctionName
func FunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}

	f := strings.Split(fn.Name(), "/")

	return f[len(f)-1]
}
