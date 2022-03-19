package services

import (
	"log"
	"testing"

	"github.com/informeai/temp-log/dto"
)

//go test -run ^TestTransform
func TestTransform(t *testing.T) {
	dtoLog := dto.Log{Type: "info", Message: "mensagem do log"}
	transfLog, err := TransformLog{}.Transform(&dtoLog)
	if err != nil {
		t.Errorf("TestTransform: got-> %s expect -> == nil\n", err.Error())
	}
	log.Println(transfLog)

}
