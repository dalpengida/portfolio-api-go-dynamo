package common

import (
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
)

// Test_Trace 함수명을 빼오게 하기 위한 기능 확인
func Test_Trace(t *testing.T) {
	file, line, functionName := Trace()
	log.Debug().
		Interface("file", file).
		Interface("line", line).
		Interface("function_name", functionName).
		Msg("success")

	f := strings.Split(functionName, "/")

	log.Debug().Interface("fn_name", f[len(f)-1]).Msg("last function name")
}

// Test_SnakeCase snake 케이스로 변경해주는 기능 확인 및 camel 로 되는지 확인
func Test_SnakeCase(t *testing.T) {
	UpperCamel := "AaBbCcDdEe000"
	lowerCamel := "a_A_b_B_c_C"

	log.Debug().Interface("upper", UpperCamel).Interface("snake", ToSnake(UpperCamel)).Msg("to snake case")
	// return : aa_bb_cc_dd_ee000

	log.Debug().Interface("lower", lowerCamel).Interface("camel", ToCamel(lowerCamel, true)).Msg("to camel case")
	// return : AABBCC
}
