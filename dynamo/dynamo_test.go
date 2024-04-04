package dynamo

import (
	"context"
	"fmt"
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
	var sliceObj []testItem
	err := dynamoClient.FindWithPK(ctx, pk, &sliceObj)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("find_obj", sliceObj).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_FindBeginsWith 는 pk, prefixSK를 이용하여 데이터를 조회 하는 기능 검사
func Test_FindBeginsWith(t *testing.T) {
	dynamoClient := New(TABLE_NAME)

	pk := "pk"
	prefixSk := "sk"
	var sliceObj []testItem
	err := dynamoClient.FindBeginsWith(ctx, pk, prefixSk, &sliceObj, 2)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("find_items", sliceObj).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_MustFindOne 는 하나의 데이터가 있을 꺼라고 믿고 조회를 시도
func Test_MustFindOne(t *testing.T) {
	dynamoClient := New(TABLE_NAME)

	pk := "pk"
	sk := "sk"

	var obj testItem
	err := dynamoClient.MustFindOne(ctx, pk, sk, &obj)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("item", obj).Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())

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

// Test_BulkPutItems 는 한번에 여러건 넣을 수 있는 기능 검사
func Test_BulkPutItems(t *testing.T) {
	items := make([]testItem, 0)
	for i := 0; i < 30; i++ {
		items = append(items, testItem{
			PK: "pk",
			SK: fmt.Sprintf("bulksk#%d", i),
		})
	}

	dynamoClient := New(TABLE_NAME)
	err := dynamoClient.PutItemsWithBatch(ctx, items)
	if err != nil {
		log.Error().Err(err).Msg("failed") // 일부러 실패를 한번 함, 초과 했을 경우를 확인 하기 위함
	}

	err = dynamoClient.PutItemsWithBatch(ctx, items[:MAX_COUNT_BULK_ITME])
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}

// Test_PutItemsWithTx 는 트랜잭션을 걸고 여러 item 을 넣을 경우 기능 검사
func Test_PutItemsWithTx(t *testing.T) {
	items := make([]testItem, 0)
	for i := 0; i < 101; i++ {
		items = append(items, testItem{
			PK: "pk",
			SK: fmt.Sprintf("bulksk#%d", i),
		})
	}

	dynamoClient := New(TABLE_NAME)
	err := dynamoClient.PutItemsWithTransaction(ctx, items)
	if err != nil {
		log.Error().Err(err).Msg("failed") // 일부러 실패를 한번 함, 초과 했을 경우를 확인 하기 위함
	}
	err = dynamoClient.PutItemsWithTransaction(ctx, items[:MAX_COUNT_TRANSACTION_ITEM])
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(TEST_SUCESS_MSG_FORMAT, common.FunctionName())
}
