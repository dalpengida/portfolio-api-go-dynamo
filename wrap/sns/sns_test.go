package sns

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/dalpengida/portfolio-go-aws/config"
	"github.com/rs/zerolog/log"
)

func Test_CreateTopic(t *testing.T) {
	err := Create(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

}

func Test_Publish(t *testing.T) {
	c := sns.NewFromConfig(config.GetAws())

	arn, _ := getTargetArn("portfolio")
	r, err := c.Publish(context.TODO(), &sns.PublishInput{
		Message: aws.String("test"),
		// Subject: , // 구독자가 email 로 구독을 했을 경우, 제목
		// PhoneNumber: , // 구독자가 sms 로 구독을 했을 경우, 수신자에 해당 하는 듯
		TargetArn: aws.String(arn),
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("response", r).Msg("success")
}
