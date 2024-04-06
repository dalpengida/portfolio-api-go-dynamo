package model

type TestItem struct {
	PK      string `dynamodbav:"pk" json:"pk"`
	SK      string `dynamodbav:"sk" json:"sk"`
	Val     string `dynamodbav:"val" json:"val"`
	Updated int64  `dynamodbav:"updated" json:"updated"`
}
