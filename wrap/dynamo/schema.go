package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CREATE_TABLE_SCHEMA 는 테이블 생성 스키마 정보
// aws 에서 recommend 하는 대로 pk, sk 만 일단 만드는 것
var CREATE_TABLE_SCHEMA = &dynamodb.CreateTableInput{

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

	// ProvisionedThroughput: &types.ProvisionedThroughput{
	// 	ReadCapacityUnits:  aws.Int64(10),
	// 	WriteCapacityUnits: aws.Int64(10),
	// },
}
