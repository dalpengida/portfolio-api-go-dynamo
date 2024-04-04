package dynamo

import (
	"context"
	"testing"
	"time"

	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/rs/zerolog/log"
)

type testItem struct {
	PK      string `dynamodbav:"pk" json:"pk"`
	SK      string `dynamodbav:"sk" json:"sk"`
	Val     string `dynamodbav:"val" json:"val"`
	Updated int64  `dynamodbav:"updated" json:"updated"`
}

var (
	TABLE_NAME             = "portfolio-test"
	ctx                    = context.Background()
	TEST_SUCESS_MSG_FORMAT = "[%s] success"
)

// Test_ListTable 는 테이블 리스트 조회 기능 테스트
func Test_ListTable(t *testing.T) {
	dynamoClient := New(TABLE_NAME)
	tables, err := dynamoClient.ListTables(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("tables", tables).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_CreateTable 는 테이블 생성 확인
func Test_CreateTable(t *testing.T) {
	dynamoClient := New(TABLE_NAME)
	tableDesc, err := dynamoClient.CreateTable(ctx, CREATE_TABLE_SCHEMA)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("table_desc", tableDesc).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// PutItem 는 아이템을 dynamo 에 upsert
func Test_PutItem(t *testing.T) {
	item := testItem{
		PK:      "pk",
		SK:      "sk",
		Val:     "val",
		Updated: time.Now().Unix(),
	}

	dynamoClient := New(TABLE_NAME)
	err := dynamoClient.PutItem(ctx, item)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_FindWithPK 는 pk 를 가지고 검색 기능 검사
func Test_FindWithPK(t *testing.T) {
	dynamoClient := New(TABLE_NAME)

	pk := "pk"
	var v []testItem
	err := dynamoClient.FindWithPK(ctx, pk, &v)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("find_obj", v).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_FindBeginsWith 는 pk, prefixSK를 이용하여 데이터를 조회 하는 기능 검사
func Test_FindBeginsWith(t *testing.T) {
	dynamoClient := New(TABLE_NAME)

	pk := "pk"
	prefixSk := "sk"
	var v []testItem
	err := dynamoClient.FindBeginsWith(ctx, pk, prefixSk, &v, 2)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("find_items", v).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_MustFindOne 는 하나의 데이터가 있을 꺼라고 믿고 조회를 시도
func Test_MustFindOne(t *testing.T) {
	dynamoClient := New(TABLE_NAME)

	pk := "pk"
	sk := "sk"

	var v testItem
	err := dynamoClient.MustFindOne(ctx, pk, sk, &v)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("item", v).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())

}

// Test_DeleteItem 는 pk, sk 를 이용하여 item 삭제
func Test_DeleteItem(t *testing.T) {
	dynamoClient := New(TABLE_NAME)

	pk := "pk"
	sk := "sk"

	err := dynamoClient.DeleteItem(ctx, pk, sk)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}
