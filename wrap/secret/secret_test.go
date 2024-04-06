package secret

import (
	"context"
	"testing"

	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/rs/zerolog/log"
)

const (
	test_secret_id          = "portfolio_id"
	test_success_msg_format = "[%s] success"
)

// Test_GetString secretmanager 에서 string으로 된 값을 조회
func Test_GetString(t *testing.T) {
	v, err := GetString(context.TODO(), test_secret_id)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("value", v).Msgf(test_success_msg_format, common.FunctionName())
}
