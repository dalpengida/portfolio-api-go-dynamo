package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/dalpengida/portfolio-go-aws/config"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	// sqs 의 경우 fifo 기능을 쓰기 위해서는 필수로 이름 끝에 fifo 가 붙어야 함
	fifo_queue_suffix = ".fifo"
)

var (
	client *sqs.Client
)

type Queue struct {
	queueName string
	queueUrl  *string
}

func init() {
	client = sqs.NewFromConfig(config.GetAws())
}

func New(queueName string) Queue {
	return Queue{
		queueName: queueName,
	}
}

// Create queue 생성
func (q Queue) Create(c context.Context, schema *sqs.CreateQueueInput) error {
	if schema == nil {
		schema = CREATE_SQS_SCHEMA
	}
	r, err := client.CreateQueue(c, schema)
	if err != nil {
		return fmt.Errorf("create queue faild, %w", err)
	}

	log.Debug().Interface("response", r).Msg("create queue success")

	return nil
}

// getUrl 는 지정한 큐의 url 정보를 조회하여 reciver 한테 저장을 해줌
// 단순하게 조회만 하는 것이 아니라 저장도 해주기 때문에 차라리 setUrl 로 함수명 변경해야 하나 고민 됨
func (q *Queue) getUrl(c context.Context) (string, error) {
	r, err := client.GetQueueUrl(c, &sqs.GetQueueUrlInput{
		QueueName: aws.String(q.queueName),
	})
	if err != nil {
		return "", fmt.Errorf("get url failed, queue : %s, %w", q.queueName, err)
	}

	log.Debug().Interface("response", r).Msg("get queue url success")
	q.queueUrl = r.QueueUrl

	return *r.QueueUrl, nil
}

//	Send 는 queue 정보를 검사를 해서 url 이 없으면 url 조회해서 넣어주고, name 을 보고 fifo 면 fifo에 맞게 전송하고 일반이면 일반으로 보내주는 함수
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

	if strings.Contains(q.queueName, fifo_queue_suffix) {
		return q.sendQueueFifo(c, string(json))
	} else {
		return q.sendQueue(c, string(json))
	}
}

// sendQueue 일반 queue 에 메시지를 전송을 해줌
func (q *Queue) sendQueue(c context.Context, message string) error {
	r, err := client.SendMessage(c, &sqs.SendMessageInput{
		QueueUrl:    q.queueUrl,
		MessageBody: aws.String(message), // message body 값의 length 가 0 이어도 오류가 남
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

// sendQueueFifo fifo queue에 맞게 메시지를 전송을 해줌
func (q *Queue) sendQueueFifo(c context.Context, message string) error {
	messageId := uuid.NewString()

	r, err := client.SendMessage(c, &sqs.SendMessageInput{
		QueueUrl:    q.queueUrl,
		MessageBody: aws.String(message), // message body 값의 length 가 0 이어도 오류가 남
		//  DelaySeconds: 0, // 0: 즉시 노출, 이외: 시간 만큼 있다가 노출
		MessageDeduplicationId: aws.String(messageId), // FIFO 타입에서는 필수, 중복 방지를 위하여 사용되는 값
		MessageGroupId:         aws.String(messageId), // FIFO 타입에서는 필수
	})
	if err != nil {
		return fmt.Errorf("fifo queue send message faeild, %w", err)
	}

	log.Debug().Interface("response", r).Msg("fifo queue send success")

	return nil
}

// BulkSend 한번에 여러 메시지를 전송을 요청을 하는 기능
// 최대 10개 까지만 가능, 그래서 내부적으로 10개가 넘을 경우 10개만 보내고 나머지는 나중에 보내는 형식으로 함
func (q *Queue) BulkSend(c context.Context, entries []types.SendMessageBatchRequestEntry) error {
	const MAX_SQS_REQUEST_COUNT = 10

	for {
		requestCount := MAX_SQS_REQUEST_COUNT
		entryCount := len(entries)
		if entryCount == 0 {
			break
		}

		var requestEnties []types.SendMessageBatchRequestEntry
		if entryCount < MAX_SQS_REQUEST_COUNT {
			requestCount = entryCount
		}

		//entries, requestEnties = common.SliceShift[types.SendMessageBatchRequestEntry](entries, requestCount)
		entries, requestEnties = common.SliceShift[types.SendMessageBatchRequestEntry](entries, requestCount)
		r, err := client.SendMessageBatch(c, &sqs.SendMessageBatchInput{
			Entries:  requestEnties,
			QueueUrl: q.queueUrl,
		})
		if err != nil {
			fmt.Println(err)
			//return err
		}
		if len(r.Failed) > 0 {
			fmt.Printf("send sqs message batch failed, %v", r.Failed)
		}
	}

	return nil
}
