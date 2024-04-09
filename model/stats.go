package model

import (
	"context"

	"github.com/dalpengida/portfolio-go-aws/config"
	"github.com/dalpengida/portfolio-go-aws/wrap/dynamo"
)

const (
	LOG_TYPE_RETENTION = "retention"
)

var ()

type Stats struct {
	TimeStamp int64  `json:"timestamp"`
	UserId    string `json:"user_id"`
	LogType   string `json:"log_type"`
	Val       string `json:"val"`
}

func (s Stats) Put(c context.Context) error {
	repo := dynamo.New(config.TABLE_LOG)
	return repo.PutItem(c, s)
}
