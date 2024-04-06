package sqs

import "github.com/aws/aws-sdk-go-v2/service/sqs"

// CREATE_SQS_SCHEMA 기본적인 sqs schema 정보를 세팅
//   - QueueName *string // 필수
//   - DelaySeconds  //  0 ~ 900 초, 지연 시간
//   - MaximumMessageSize // An integer from 1,024 bytes (1 KiB) to 262,144 bytes (256 KiB). Default: 262,144 (256 KiB).
//   - MessageRetentionPeriod // An integer from 60 seconds (1 minute) to 1,209,600 seconds (14 days). Default: 345,600 (4 days).
//   - Policy // Policies (https://docs.aws.amazon.com/IAM/latest/UserGuide/PoliciesOverview.html)
//   - ReceiveMessageWaitTimeSeconds // An integer from 0 to 20 (seconds). Default: 0.
//   - VisibilityTimeout – An integer from 0 to 43,200 (12 hours). Default: 30
//     https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html
//   - RedrivePolicy //dlq 관련 정책
//   - deadLetterTargetArn // dlq 관련 arn 정보
//   - maxReceiveCount // 실패 했을 때 몇번까지 재 시도 하고 dlq 로 넘어갈 것인가 관련 설정, Default: 10.

var CREATE_SQS_SCHEMA = &sqs.CreateQueueInput{
	Attributes: map[string]string{
		"DelaySeconds":           "0",
		"MessageRetentionPeriod": "86400",
		"VisibilityTimeout":      "0",
	},
}
