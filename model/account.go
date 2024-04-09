package model

import (
	"context"
	"time"

	"github.com/dalpengida/portfolio-go-aws/wrap/dynamo"
	"github.com/google/uuid"
)

const (
	//	prefix_account_pk = "user#"
	prefix_account_sk = "account#"
)

var ()

type Account struct {
	// 클라에서 사용을 하기 위하여 json 파싱시 사용, pk 값이 변경이 되면 꼭 확인을 해야 함
	PK        string `dynamodbav:"pk" json:"-"`
	SK        string `dynamodbav:"sk" json:"-"`
	UserId    string `dynamodbav:"user_id" json:"user_id"`
	LastLogin int64  `dynamodbav:"last_login" json:"last_login"`
	Created   int64  `dynamodbav:"exp" json:"exp"`
	Updated   int64  `dynamodbav:"updated" json:"-"`
}

func NewAccount() Account {
	userId := uuid.NewString()
	now := time.Now().Unix()

	return Account{
		PK:      userId,
		SK:      prefix_account_sk,
		UserId:  userId,
		Created: now,
		Updated: now,
	}
}

// SKForAccount account 의 sk 값
func SKForAccount() string {
	return prefix_account_sk
}

// Put item 데이터를 upsert
func (a *Account) Put(c context.Context, item Account) error {
	repo := dynamo.NewDefault()

	return repo.PutItem(c, item)
}

// Remove 유저 정보에 맞는 데이터 삭제
func (a *Account) Remove(c context.Context, item Account) error {
	repo := dynamo.NewDefault()

	return repo.DeleteItem(c, item.PK, item.SK)
}

// Find pk, sk 를 이용해서 account 정보 조회
func (a Account) Find(c context.Context, userId string) (Account, error) {
	var r Account
	repo := dynamo.NewDefault()
	err := repo.MustFindOne(c, userId, prefix_account_sk, &r)

	return r, err
}
