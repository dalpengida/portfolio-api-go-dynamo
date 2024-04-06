package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	awsConfig aws.Config
)

func init() {
	var err error
	// 한번만 읽어 주면 되기 떄문에 초기화 할 때 하는 걸로 함
	if !common.IsAWSLambda() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Error().Msg("load .env file failed")
		}
	}

	// 각 서비스별로 매번 하지 말고, 여기서 초기화 할떄 값을 가져와서 다른데서는 필요할 때 이 값을 가져가서 사용하게 함
	awsConfig, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err) // config 못 읽으면 서비스 장애임
	}
}

// Config 는 환경변수에서 값을 읽어줌
// .env 파일이 있을 때는 해당 값을 환경변수에서 읽은 거 처럼 해줌
func Config(key string) string {
	return os.Getenv(key)
}

// GetAws awsConfig 값을 전달
// 어차피 한번만 호출해놓으면 모든 aws service들을 초기화 할떄 동일한 값으로 사용하기 떄문에
// 머 크게 차이가 나지 않겠지만 한번만 호출해서 땡겨가서 쓰도록 함
func GetAws() aws.Config {
	return awsConfig
}
