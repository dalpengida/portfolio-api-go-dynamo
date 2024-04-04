package common

import (
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
)

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
