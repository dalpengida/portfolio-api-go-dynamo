package sqs

import (
	"context"
	"testing"

	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/dalpengida/portfolio-go-aws/model"
	"github.com/rs/zerolog/log"
)

var (
	TEST_QUEUE_NAME        = "portfolio"
	TEST_QUEUE_FIFO_NAME   = "portfolio.fifo"
	TEST_SUCESS_MSG_FORMAT = "[%s] success"
)

// Test_GetUrl 는 queue 에 요청을 보낼때 사용될 url 정보를 조회 기능 검사
func Test_GetUrl(t *testing.T) {
	queue := New(TEST_QUEUE_NAME)
	url, err := queue.getUrl(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("url", url).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_Send 는 기본적인 sqs에 message 를 전송하는 기능 검사
func Test_Send(t *testing.T) {
	queue := New(TEST_QUEUE_NAME)

	var item model.TestItem
	item.PK = "pk"
	item.SK = "sk"
	item.Val = "val"

	err := queue.Send(context.TODO(), item)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_SendToFifoQueue 는 fifo queue에 message 를 전송하는 기능 검사
func Test_SendToFifoQueue(t *testing.T) {
	queue := New(TEST_QUEUE_FIFO_NAME)

	var item model.TestItem
	item.PK = "pk"
	item.SK = "sk"
	item.Val = "val"

	err := queue.Send(context.TODO(), item)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}
