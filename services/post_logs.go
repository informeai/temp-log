package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/informeai/temp-log/dto"
)

//TransformLog is struct for transform logs
type TransformLog struct {
	Message string
}

//NewTransformLog return new instance of TransformLog
func NewTransformLog() *TransformLog {
	return &TransformLog{}
}

//Transform is return log transformed
func (t TransformLog) Transform(log *dto.Log) (TransformLog, error) {
	if len(log.Type) < 3 || len(log.Message) < 5 {
		return TransformLog{}, errors.New("length of type or message not permited")
	}
	if len(log.Project) > 1 {
		t.Message = fmt.Sprintf("[%v]:[%s]:[%s] %s", time.Now().UTC().Format(time.RFC3339), strings.ToUpper(log.Project), strings.ToUpper(log.Type), log.Message)
		return t, nil
	}
	t.Message = fmt.Sprintf("[%v]:[DEFAULT]:[%s] %s", time.Now().UTC().Format(time.RFC3339), strings.ToUpper(log.Type), log.Message)
	return t, nil
}
