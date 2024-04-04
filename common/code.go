package common

import "errors"

// 공통으로 쓰는 에러 코드들을 정리
// TODO: code 라는 패키지로 뺄까 고민 중
var (
	ErrorNotFountItem           = errors.New("not found item")
	ErrorRequestParameterExceed = errors.New("request parameter exceed")
)
