package sns

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/dalpengida/portfolio-go-aws/config"
	"github.com/rs/zerolog/log"
)

const (
	seperator = ":"
)

var (
	client *sns.Client
	topics map[string]string
)

type Notification struct {
	topic     string
	targetArn string
}

func init() {
	topics = make(map[string]string, 0)

	client = sns.NewFromConfig(config.GetAws())
	err := topicsToMap()
	if err != nil {
		panic(err)
	}
}

func New(topic string) Notification {
	arn, ok := topics[topic]
	if !ok {
		// aws topic 리스트에 없는 애를 호출을 하면 잘 못 들고 온거라고 판단을 함
		// 잘 못된 topic 을 가져온 거라 panic
		panic(fmt.Errorf("invlid topic, [%s] is not found topic list", topic))
	}

	return Notification{
		topic:     topic,
		targetArn: arn,
	}
}

// Create topic 생성 함수
// https://docs.aws.amazon.com/sns/latest/dg/sns-create-topic.html
func Create(c context.Context) error {
	r, err := client.CreateTopic(c, &sns.CreateTopicInput{
		Name: aws.String("portfolio_test"),
	})
	if err != nil {
		return fmt.Errorf("create sns topic failed, %w", err)
	}

	log.Debug().Interface("response", r).Msg("sns topice create success")

	return nil
}

// topicsToMap targetArn 을 가져오기 위하여 aws sns topic 정보들을 모두 가져와서 map으로 가지고 있음
func topicsToMap() error {
	r, err := client.ListTopics(context.TODO(), &sns.ListTopicsInput{})
	if err != nil {
		return fmt.Errorf("get list topics failed, %w", err)
	}

	log.Debug().Interface("response", r).Msg("get list topics success")

	for _, v := range r.Topics {
		sp := strings.Split(*v.TopicArn, seperator)
		topics[sp[len(sp)-1]] = *v.TopicArn
	}

	log.Debug().Interface("topics", topics).Msg("")

	return nil
}

// getTargetArn 는 미리 만들어 놓은 topic map 에서 topic arn 을 찾아서 넘겨 줌
func getTargetArn(topic string) (v string, ok bool) {
	v, ok = topics[topic], false
	return
}

// Publish 는 sns 로 메시지 전달 , 발행, 전송
// Subject: , // 구독자가 email 로 구독을 했을 경우, 제목
// PhoneNumber: , // 구독자가 sms 로 구독을 했을 경우, 수신자에 해당 하는 듯
func (n Notification) Publish(c context.Context, message string) error {
	r, err := client.Publish(c, &sns.PublishInput{
		Message:   aws.String(message),
		TargetArn: aws.String(n.targetArn),
	})
	if err != nil {
		return fmt.Errorf("sns publish failed, %w", err)
	}

	log.Debug().Interface("response", r).Msg("sns publish success")

	return nil
}

// SubscribeTopic 지정 topic 에 구독을 신청을 함
func (n Notification) SubscribeTopic(c context.Context, protocol, endpoint string) error {
	if !isValidSubscribeProtocol(protocol) {
		return fmt.Errorf("invalid protocol, %s", protocol)
	}

	r, err := client.Subscribe(c, &sns.SubscribeInput{
		// http, https, email, email-json, sms, sqs, application, lambda, firehouse
		TopicArn:              aws.String(n.targetArn),
		Protocol:              aws.String(protocol),
		Endpoint:              aws.String(endpoint), // protocol 에 따라서 구독을 받을 대상 정보
		ReturnSubscriptionArn: true,
	})
	if err != nil {
		return fmt.Errorf("topic subscribe failed, topic : %s , %w", n.topic, err)
	}

	log.Debug().Interface("response", r).Msg("subscribe success")

	return nil
}

// UnsubscribeTopic 는 구독한 arn 를 가지고 구독 해제를 함
func UnsubscribeTopic(c context.Context, subscribeArn string) error {
	r, err := client.Unsubscribe(c, &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(subscribeArn),
	})
	if err != nil {
		return fmt.Errorf("unsubscribe failed, arn : %s, %w", subscribeArn, err)
	}

	log.Debug().Interface("response", r).Msg("unsubsribe success")

	return nil
}

// isValidSubscribeProtocol topic 구독 가능한 프로토콜인지 확인
func isValidSubscribeProtocol(protocol string) bool {
	switch protocol {
	case "http", "https", "email", "email-json", "sms", "sqs", "application", "lambda", "firehouse":
		return true
	default:
		return false
	}
}
