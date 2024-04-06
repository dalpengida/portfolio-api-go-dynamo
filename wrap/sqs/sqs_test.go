package sqs

import (
	"context"
	"testing"

	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/dalpengida/portfolio-go-aws/model"
	"github.com/rs/zerolog/log"
)

const (
	test_queue_name         = "portfolio"
	test_queue_fifo_name    = "portfolio.fifo"
	test_success_msg_format = "[%s] success"
)

// Test_GetUrl 는 queue 에 요청을 보낼때 사용될 url 정보를 조회 기능 검사
func Test_GetUrl(t *testing.T) {
	queue := New(test_queue_name)
	url, err := queue.getUrl(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("url", url).Msgf(test_success_msg_format, common.FunctionName())
}

// Test_Send 는 기본적인 sqs에 message 를 전송하는 기능 검사
func Test_Send(t *testing.T) {
	queue := New(test_queue_name)

	var item model.TestItem
	item.PK = "pk"
	item.SK = "sk"
	item.Val = "val"

	err := queue.Send(context.TODO(), item)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(test_success_msg_format, common.FunctionName())
}

// Test_SendToFifoQueue 는 fifo queue에 message 를 전송하는 기능 검사
func Test_SendToFifoQueue(t *testing.T) {
	queue := New(test_queue_fifo_name)

	var item model.TestItem
	item.PK = "pk"
	item.SK = "sk"
	item.Val = "val"

	err := queue.Send(context.TODO(), item)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(test_success_msg_format, common.FunctionName())
}
