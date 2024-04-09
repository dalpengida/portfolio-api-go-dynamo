package config

import (
	"fmt"
)

const (
	TABLE_LOG = "portfolio-log"
	STAGE     = "STAGE"
)

// AccountTopicName account topic 이름을 전달
func AccountTopicName() string {
	return fmt.Sprintf("topic-%s-account", Config(STAGE))
}
