package service

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
)

type GenPass interface {
	GeneratePassword(length int) (string, error)
}

type genPassService struct {
}

func NewGenPass() GenPass {
	return &genPassService{}
}

func (g *genPassService) GeneratePassword(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}
	dash := "-"
	result := strings.ReplaceAll(uuid.New().String(), dash, "")
	pass_len := len(result)
	for (pass_len < length) {
		result += strings.ReplaceAll(uuid.New().String(), dash, "")
		pass_len = len(result)
	}
	result = result[:length]
	return result, nil
}
