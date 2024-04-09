package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/dalpengida/portfolio-go-aws/model"
)

func init() {
	lambda.Start(handler)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		var noti model.AccountNoti
		err := json.Unmarshal([]byte(record.Body), &noti)
		if err != nil {
			return err
		}
		// 같은 날 들어온 데이터라고 하면 그냥 넘김
		if !common.IsDiffDate(noti.PreLastLogin, noti.LastLogin) {
			continue
		}

		// 날짜가 다르면 retention 로그를 일단 하나 남김
		stats := model.Stats{
			TimeStamp: time.Now().Unix(),
			UserId:    noti.UserId,
			LogType:   model.LOG_TYPE_RETENTION,
			Val:       record.Body,
		}
		err = stats.Put(context.TODO())
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
}
