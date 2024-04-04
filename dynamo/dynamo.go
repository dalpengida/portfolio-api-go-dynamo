package dynamo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/dalpengida/portfolio-go-aws/common"
	"github.com/rs/zerolog/log"
)

var (
	awsConfig aws.Config
	client    *dynamodb.Client
)

// TableBasics 는 dynamodb wrapping 한 기능들을 사용을 할 때, 다른 테이블을 실수로 사용하게 되는 것을 방지 하기 위함
// 여러 테이블을 사용할 수 있을 것 같아서 table 값을 초기화 할때 받아서 사용하게 하기 위함
type TableBasics struct {
	tableName string
}

func init() {
	var err error
	awsConfig, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client = dynamodb.NewFromConfig(awsConfig)
}

func New(tablename string) TableBasics {
	return TableBasics{
		tableName: tablename,
	}
}

// CreateTable 는 테이블 생성을 해주는 함수
// dynamo 는 pk, sk 를 제외하고는 언제든지 attribute 가 변경 될 수 있어서 그냥 pk, sk 만 대충 생성을 해도 되는 듯
// 스키마도 크게 변경될 일이 없을 것 같지만, 혹시 몰라서 받는 걸로
func (t TableBasics) CreateTable(c context.Context, createTableSchema *dynamodb.CreateTableInput) (*types.TableDescription, error) {
	if createTableSchema == nil {
		// 임시 테이블 스키마
		createTableSchema = CREATE_TABLE_SCHEMA

	}
	createTableSchema.TableName = aws.String(t.tableName)

	r, err := client.CreateTable(c, createTableSchema)
	if err != nil {
		return nil, fmt.Errorf("create table %v failed, %w", t.tableName, err)

	} else {
		waiter := dynamodb.NewTableExistsWaiter(client)
		// 생성을 요청하고 바로 올라오는게 아니다 보니 좀 대기
		err = waiter.Wait(c, &dynamodb.DescribeTableInput{
			TableName: aws.String(t.tableName)}, 5*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("wait for table exists failed, %w", err)
		}
	}

	log.Debug().Interface("r", r).Msg("create table response")

	return r.TableDescription, err
}

// IsExist 테이블 존재 여부 확인
func (t TableBasics) IsExist(c context.Context) (bool, error) {
	_, err := client.DescribeTable(
		c, &dynamodb.DescribeTableInput{TableName: aws.String(t.tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			return false, fmt.Errorf("table %v does not exist, %w", t.tableName, err)

		} else {
			return false, fmt.Errorf("couldn't determine existence of table %v", t.tableName)
		}
	}

	return true, nil
}

// ListTables 테이블 리스트 조회
func (TableBasics) ListTables(c context.Context) ([]string, error) {
	r, err := client.ListTables(c, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("table list lookup failed, %w", err)
	}

	log.Debug().Interface("tables", r.TableNames).Msg("")
	return r.TableNames, nil
}

// PutItem 는 item interface를 받아서 데이터를 추가
// dynamo 에서 putitem 은 upsert 인것으로 확인
// response 값은 쓸일이 없을 것 같아서 생략
func (t TableBasics) PutItem(c context.Context, item interface{}) error {
	i, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("attribute marshal map failed, %w", err)
	}
	response, err := client.PutItem(c, &dynamodb.PutItemInput{
		TableName: aws.String(t.tableName), Item: i,
	})
	if err != nil {
		return fmt.Errorf("put item failed, %w", err)
	}

	log.Debug().Interface("response", response).Msg("put item success")

	return nil
}

// FindWithPK 는 pk 를 기준으로 데이터를 모두 조회 , 입력한 구조체로 바인딩을 해서 전달
// pk 를 기준으로 조회를 하다 보면 메시지는 여러건이 나오기 때문에 slice obj 형태로 인자를 받아야 함
func (t TableBasics) FindWithPK(c context.Context, pk string, sliceObj interface{}) error {
	// TODO: sliceObj 검사 로직을 넣어야 함
	// 내부 UnmarshalListOfMaps 에서 걸러질 수 있지만, 고민

	response, err := client.Query(c, &dynamodb.QueryInput{
		TableName: aws.String(t.tableName),

		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("find with pk failed, %w", err)
	}

	// response.items 는 []map[string]types.AttributeValue
	// 외부에서 바로 사용을 할 수 있도록 binding
	err = attributevalue.UnmarshalListOfMaps(response.Items, sliceObj)
	if err != nil {
		return fmt.Errorf("attributevalue unmarshallistofmaps failed, err : %w", err)
	}

	log.Debug().Interface("pk", pk).Interface("response", response).Msg("find with pk success")

	return nil
}

// FindBeginsWith 는 pk, sk의 prefix 값을 이용하여 검색을 하는 로직
// object interface{}의 경우는 특정 struct만 들어갈 수 있도록 하던가, 아니면 interface로 만들어서 방어해야함
// 일단은 slice interface 형태로 사용해야함
// 'begins_with' function 은 대소문자 구분함, 괜히 예약어라고 해서 upper case 로 섰다가 망함
// limit 값으로 한건만 찾아야 하는 경우, 그리고 여러건을 찾아야 하는 경우를 함수를 나눠서 사용할까 했지만 어차피 binding 할 때 slice 로 돌려 주기 때문에 의미 없음
func (t TableBasics) FindBeginsWith(c context.Context, pk, prefixSk string, objSlice interface{}, limit int) error {
	response, err := client.Query(c, &dynamodb.QueryInput{
		TableName: aws.String(t.tableName),

		KeyConditionExpression: aws.String("pk = :pk and begins_with(sk, :beginsWith)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
			":beginsWith": &types.AttributeValueMemberS{
				Value: prefixSk,
			},
		},
		Limit: aws.Int32(int32(limit)),
	})
	if err != nil {
		return fmt.Errorf("find beginswith failed, %w", err)
	}

	// response.items 는 []map[string]types.AttributeValue 임
	err = attributevalue.UnmarshalListOfMaps(response.Items, objSlice)
	if err != nil {
		return fmt.Errorf("attributevalue unmarshallistofmaps failed, err : %w", err)
	}

	log.Debug().Interface("pk", pk).Interface("prefix_sk", prefixSk).Interface("response", response).Msg("find begins with success")

	return nil
}

// MustFindOne 는 pk, sk를 이용하여 한 데이터만 찾기 위한 함수, 지정한 struct 구조로 바인딩하여 전달
// 하나라도 없으면 걍 에러 처리 왜냐? must 이기 떄문
func (t TableBasics) MustFindOne(c context.Context, pk, sk string, obj interface{}) error {
	// sk 값이 따로 없을 경우, #으로 지정을 해서 조회
	if sk == "" {
		sk = "#"
	}

	response, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(t.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("couldn`t get info about pk : %v, sk : %v, err : %w", pk, sk, err)
	}
	if response.Item == nil {
		log.Error().Interface("pk", pk).Interface("sk", sk).Msg("not found item")
		return common.ErrorNotFountItem
	}

	err = attributevalue.UnmarshalMap(response.Item, obj)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal response, err : %w", err)

	}

	return nil
}

// DeleteItem 는 pk, sk 를 인자로 받아서 item 삭제
func (t TableBasics) DeleteItem(c context.Context, pk, sk string) error {
	_, err := client.DeleteItem(c, &dynamodb.DeleteItemInput{
		TableName: aws.String(t.tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("item delete failed, %w", err)
	}

	log.Debug().Interface("pk", pk).Interface("sk", sk).Msg("delete item success")

	return nil
}
