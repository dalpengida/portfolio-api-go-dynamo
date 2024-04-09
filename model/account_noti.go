package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dalpengida/portfolio-go-aws/config"
	"github.com/dalpengida/portfolio-go-aws/wrap/sns"
)

const ()

var ()

type AccountNoti struct {
	UserId       string `json:"user_id"`
	PreLastLogin int64  `json:"pre_last_login"`
	LastLogin    int64  `json:"last_login"`
	EventType    string `json:"event_type"`
	TimeStamp    int64  `json:"timestamp"`
}

// func (a *AccountNoti) Bind(old, new Account) {
// 	a.AccountBase = new.AccountBase
// 	a.PreLastLogin = old.LastLogin
// }

func NewAccountNoti() AccountNoti {
	return AccountNoti{}
}

// Publish 는 AccountNoti 구조체의 데이터를 json 으로 marshaling 해서 sns publish 함
func (a AccountNoti) Publish(c context.Context) error {
	notiMessage, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("account noti publish failed, %w", err)
	}

	topic := sns.New(config.AccountTopicName())

	return topic.Publish(c, string(notiMessage))
}
