package common

import (
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
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

// ToSnake 는 snake case 로 변환을 시켜 줌
func ToSnake(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1}_${2}")
	}
	return strings.ToLower(s)
}

// ToCamel 는 camel case 형태로 변환을 시켜 줌
func ToCamel(s string, uppercaseFirstLetter bool) string {
	if len(s) == 0 {
		return s
	}
	replFunc := func(w string) string {
		if strings.HasPrefix(w, "_") && !strings.HasPrefix(w, "__") {
			return strings.ToUpper(w[1:])
		}
		return strings.ToUpper(w)
	}

	if uppercaseFirstLetter {
		re := regexp.MustCompile(`(?:^|_)(.)`)
		return re.ReplaceAllStringFunc(s, replFunc)
	} else {
		return strings.ToLower(string(s[0])) + ToCamel(s, true)[1:]
	}
}

// SlicePop 는 slice 형태를 가지고 있는 그 어떤 형태에서도 pop을 해주고 남은 slice 를 리턴
func SlicePop[T any](orig []T, i int) ([]T, T) {
	elem := orig[i]
	orig = append(orig[:i], orig[i+1:]...)
	return orig, elem
}

// SlicePopSlice 는 slice 형태를 가지고 있는 그 어떤 형태에서도 slice 값을 pop을 해주고 남은 slice 를 리턴
func SlicePopSlice[T any](orig []T, start, end int) ([]T, []T) {
	elem := append(make([]T, 0, len(orig)), orig[start:end]...)
	orig = append(orig[:start], orig[end:]...)
	return orig, elem
}

// SliceShift 는 slice 형태를 가지고 있는 구조체를 i 만큼 시프트 시켜줌
func SliceShift[T any](s []T, i int) ([]T, []T) {
	return s[i:], s[:i]
}

// SliceCopy 는 slice 들의 깊은 복사를 해주는 것
// 속도 관련해서는 util_test.go 에 benchmark 돌린게 있음
func SliceCopy[T any](orig []T) []T {
	return append(make([]T, 0, len(orig)), orig...)
}

// IsAWSLambda 환경이 aws lambda 인지 확인
func IsAWSLambda() bool {
	return os.Getenv("AWS_EXECUTION_ENV") == "AWS_Lambda_go1.x"
}

// IsDiffDate 는 unixTimestamp를 받아서 같은 날짜인지 확인
func IsDiffDate(t1, t2 int64) bool {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)

	if tt1.Year() != tt2.Year() || tt1.Month() != tt2.Month() || tt1.Day() != tt2.Day() {
		return true
	}

	return false
}
