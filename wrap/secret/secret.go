package secret

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/dalpengida/portfolio-go-aws/config"
	"github.com/rs/zerolog/log"
)

const (
	latest_version_for_aws_secretsmanager = "AWSCURRENT"
)

var (
	client *secretsmanager.Client
)

func init() {
	client = secretsmanager.NewFromConfig(config.GetAws())
}

// GetString secretmanager 에서 값을 가져옴
func GetString(c context.Context, secretId string) (string, error) {
	r, err := client.GetSecretValue(c, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretId),
		VersionStage: aws.String(latest_version_for_aws_secretsmanager),
	})
	if err != nil {
		return "", fmt.Errorf("get secretmanager value failed, secret id : %s, %w", secretId, err)
	}

	log.Debug().Interface("response", r).Msgf("secret manager get value success, secret id : %s", secretId)

	return *r.SecretString, nil
}

// GetWithJsonUnmarshal secretmanger에서 값을 가져옴
// 들어가 있는 값이 json marshaling 이 되었을 경우 여기서 unmarshaling 해서 던져 줌
func GetWithJsonUnmarshal(c context.Context, secretId string, obj interface{}) error {
	r, err := GetString(c, secretId)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(r), obj)
	if err != nil {
		return fmt.Errorf("data unmarshaling failed, %w", err)
	}

	return nil
}
