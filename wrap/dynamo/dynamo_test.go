package dynamo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/dalpengida/portfolio-go-aws/model"
	"github.com/rs/zerolog/log"
)

var (
	test_table_name         = "portfolio-test"
	test_success_msg_format = "[%s] success"
)

// Test_ListTable 는 테이블 리스트 조회 기능 테스트
func Test_ListTable(t *testing.T) {
	dynamoClient := New(test_table_name)
	tables, err := dynamoClient.ListTables(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("tables", tables).Msgf(test_success_msg_format, common.FunctionName())
}

// Test_CreateTable 는 테이블 생성 확인
func Test_CreateTable(t *testing.T) {
	dynamoClient := New(test_table_name)
	tableDesc, err := dynamoClient.CreateTable(context.TODO(), CREATE_TABLE_SCHEMA)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("table_desc", tableDesc).Msgf(test_success_msg_format, common.FunctionName())
}

// PutItem 는 아이템을 dynamo 에 upsert
func Test_PutItem(t *testing.T) {
	item := model.TestItem{
		PK:      "pk",
		SK:      "sk",
		Val:     "val",
		Updated: time.Now().Unix(),
	}

	dynamoClient := New(test_table_name)
	err := dynamoClient.PutItem(context.TODO(), item)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(test_success_msg_format, common.FunctionName())
}

// Test_FindWithPK 는 pk 를 가지고 검색 기능 검사
func Test_FindWithPK(t *testing.T) {
	dynamoClient := New(test_table_name)

	pk := "pk"
	var sliceObj []model.TestItem
	err := dynamoClient.FindWithPK(context.TODO(), pk, &sliceObj)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("find_obj", sliceObj).Msgf(test_success_msg_format, common.FunctionName())
}

// Test_FindBeginsWith 는 pk, prefixSK를 이용하여 데이터를 조회 하는 기능 검사
func Test_FindBeginsWith(t *testing.T) {
	dynamoClient := New(test_table_name)

	pk := "pk"
	prefixSk := "sk"
	var sliceObj []model.TestItem
	err := dynamoClient.FindBeginsWith(context.TODO(), pk, prefixSk, &sliceObj, 2)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("find_items", sliceObj).Msgf(test_success_msg_format, common.FunctionName())
}

// Test_MustFindOne 는 하나의 데이터가 있을 꺼라고 믿고 조회를 시도
func Test_MustFindOne(t *testing.T) {
	dynamoClient := New(test_table_name)

	pk := "pk"
	sk := "sk"

	var obj model.TestItem
	err := dynamoClient.MustFindOne(context.TODO(), pk, sk, &obj)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Interface("item", obj).Msgf(test_success_msg_format, common.FunctionName())

}

// Test_DeleteItem 는 pk, sk 를 이용하여 item 삭제
func Test_DeleteItem(t *testing.T) {
	dynamoClient := New(test_table_name)

	pk := "pk"
	sk := "sk"

	err := dynamoClient.DeleteItem(context.TODO(), pk, sk)
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(test_success_msg_format, common.FunctionName())
}

// Test_BulkPutItems 는 한번에 여러건 넣을 수 있는 기능 검사
func Test_BulkPutItems(t *testing.T) {
	items := make([]model.TestItem, 0)
	for i := 0; i < 30; i++ {
		items = append(items, model.TestItem{
			PK: "pk",
			SK: fmt.Sprintf("bulksk#%d", i),
		})
	}

	dynamoClient := New(test_table_name)
	err := dynamoClient.PutItemsWithBatch(context.TODO(), items)
	if err != nil {
		log.Error().Err(err).Msg("failed") // 일부러 실패를 한번 함, 초과 했을 경우를 확인 하기 위함
	}

	err = dynamoClient.PutItemsWithBatch(context.TODO(), items[:max_count_bulk_item])
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(test_success_msg_format, common.FunctionName())
}

// Test_PutItemsWithTx 는 트랜잭션을 걸고 여러 item 을 넣을 경우 기능 검사
func Test_PutItemsWithTx(t *testing.T) {
	items := make([]model.TestItem, 0)
	for i := 0; i < 101; i++ {
		items = append(items, model.TestItem{
			PK: "pk",
			SK: fmt.Sprintf("bulksk#%d", i),
		})
	}

	dynamoClient := New(test_table_name)
	err := dynamoClient.PutItemsWithTransaction(context.TODO(), items)
	if err != nil {
		log.Error().Err(err).Msg("failed") // 일부러 실패를 한번 함, 초과 했을 경우를 확인 하기 위함
	}
	err = dynamoClient.PutItemsWithTransaction(context.TODO(), items[:max_count_transaction_item])
	if err != nil {
		t.Fatal(err)
	}

	log.Debug().Msgf(test_success_msg_format, common.FunctionName())
}
