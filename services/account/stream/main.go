package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"

	"github.com/dalpengida/portfolio-go-aws/model"
)

func handler(ctx context.Context, event events.DynamoDBEvent) error {
	for _, record := range event.Records {
		sk := record.Change.Keys["sk"].String()

		if sk != model.SKForAccount() {
			continue
		}

		// 유저 진입 알림 및 last login 계산을 위한 raw 데이터
		switch record.EventName {
		case "INSERT", "MODIFY":
			preLastLogin, err := record.Change.OldImage["last_login"].Int64()
			if err != nil {
				return err
			}
			lastLogin, err := record.Change.NewImage["last_login"].Int64()
			if err != nil {
				return err
			}

			notiMessage := model.AccountNoti{
				UserId:       record.Change.NewImage["user_id"].String(),
				PreLastLogin: preLastLogin,
				LastLogin:    lastLogin,
				EventType:    record.EventName,
				TimeStamp:    time.Now().Unix(),
			}
			err = notiMessage.Publish(context.TODO())
			if err != nil {
				return err
			}

			log.Debug().Interface("noti_message", notiMessage).Msg("noti message publish success")

		case "REMOVE":
		}
	}

	return nil
}

func init() {
	lambda.Start(handler)
}

func main() {
}
