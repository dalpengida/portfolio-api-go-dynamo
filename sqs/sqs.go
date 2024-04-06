package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rs/zerolog/log"
)

var (
	awsConfig aws.Config
	client    *sqs.Client
)

type Queue struct {
	QueueName string
	queueUrl  *string
}

func init() {
	var err error
	awsConfig, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client = sqs.NewFromConfig(awsConfig)
}

func New(queueName string) Queue {
	return Queue{
		QueueName: queueName,
	}
}

// getUrl 는 지정한 큐의 url 정보를 조회하여 reciver 한테 저장을 해줌
// 단순하게 조회만 하는 것이 아니라 저장도 해주기 때문에 차라리 setUrl 로 함수명 변경해야 하나 고민 됨
func (q *Queue) getUrl(c context.Context) (string, error) {
	r, err := client.GetQueueUrl(c, &sqs.GetQueueUrlInput{
		QueueName: aws.String(q.QueueName),
	})
	if err != nil {
		return "", fmt.Errorf("get url failed, queue : %s, %w", q.QueueName, err)
	}

	log.Debug().Interface("response", r).Msg("get queue url success")
	q.queueUrl = r.QueueUrl

	return *r.QueueUrl, nil
}

//	Send
//
// 내부를 확인을 해보면 전송 타입에 따른 정보가 있지만, 아직은 v2에서는 제대로 구현이 되어 있지 않은 것으로 보임
// 그래서 data 구조체를 넘겨야 할 경우 json 으로 marshaling 해서 전달을 하고 받는 쪽에서 다시 unmarshaling 하는 것으로
// // This member is required.
// // DataType *string,BinaryListValues [][]byte,	BinaryValue []byte, StringListValues []string
func (q *Queue) Send(c context.Context, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("invalid obj or obj is nil")
	}

	if q.queueUrl == nil {
		_, err := q.getUrl(c)
		if err != nil {
			return err
		}
	}

	json, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("queue message json marshaling faeild, %w", err)
	}

	r, err := client.SendMessage(c, &sqs.SendMessageInput{
		QueueUrl:    q.queueUrl,
		MessageBody: aws.String(string(json)), // message body 값의 length 가 0 이어도 오류가 남
		//  DelaySeconds: 0, // 0: 즉시 노출, 이외: 시간 만큼 있다가 노출
		//MessageDeduplicationId: aws.String(""), // FIFO 타입에서는 필수, 중복 방지를 위하여 사용되는 값
		// MessageGroupId:         aws.String(""), // FIFO 타입에서는 필수
	})
	if err != nil {
		return fmt.Errorf("queue send message faeild, %w", err)
	}

	log.Debug().Interface("response", r).Msg("queue send success")

	return nil
}
