package dynamo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

type testItem struct {
	PK      string `dynamodbav:"pk" json:"pk"`
	SK      string `dynamodbav:"sk" json:"sk"`
	Val     string `dynamodbav:"val" json:"val"`
	Updated int64  `dynamodbav:"updated" json:"updated"`
}

var (
	TABLE_NAME = "portfolio"
	ctx        = context.Background()

	awsConfig aws.Config
	client    *dynamodb.Client
)

func init() {
	var err error
	awsConfig, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client = dynamodb.NewFromConfig(awsConfig)
}

// Test_ListTable 는 테이블 리스트 조회 기능 테스트
func Test_ListTable(t *testing.T) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	client = dynamodb.NewFromConfig(awsConfig)
	list, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(list)
}

// Test_CreateTable 는 테이블 생성 확인
func Test_CreateTable(t *testing.T) {

	// 테이블 생성 스키마
	createTableSchema := &dynamodb.CreateTableInput{

		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("pk"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("sk"),
			AttributeType: types.ScalarAttributeTypeS,
		}},

		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("pk"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("sk"),
			KeyType:       types.KeyTypeRange,
		}},
		// on demand
		BillingMode: types.BillingModePayPerRequest,
	}

	createTableSchema.TableName = aws.String(TABLE_NAME)

	table, err := client.CreateTable(ctx, createTableSchema)
	if err != nil {
		log.Error().Err(err).Msgf("create table %v failed", TABLE_NAME)
		t.Fatal(err)

	} else {
		waiter := dynamodb.NewTableExistsWaiter(client)
		err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(TABLE_NAME)}, 5*time.Minute)

		if err != nil {
			log.Error().Err(err).Msg("wait for table exists failed")
			t.Fatal(err)
		}
	}

	log.Info().Interface("table desc", table.TableDescription).Msg("success")
}

// PutItem 는 아이템을 dynamo 에 upsert
func Test_PutItem(t *testing.T) {
	item := testItem{
		PK:      "pk",
		SK:      "sk",
		Val:     "val",
		Updated: time.Now().Unix(),
	}

	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	r, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME), Item: i,
	})
	if err != nil {
		log.Error().Err(err).Msg("put item failed")
	}

	log.Info().Interface("r", r).Msg("success")
}

// Test_Find 는 pk 를 가지고 검색 기능 검사
func Test_Find(t *testing.T) {
	r, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(TABLE_NAME),

		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: "pk",
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("find failed")
	}

	log.Info().Interface("r", r).Msg("success")
}
