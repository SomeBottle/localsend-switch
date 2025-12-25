package utils

// 发现包相关的工具函数

import (
	"crypto/rand"
)

const ID_LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomSwitchID 生成一个随机的 16 字节 Switch ID，用于标识交换节点
func GenerateRandomSwitchID() string {
	randomBytes := make([]byte, 16)
	_, _ = rand.Read(randomBytes)
	for i, b := range randomBytes {
		randomBytes[i] = ID_LETTERS[int(b)%len(ID_LETTERS)]
	}
	return string(randomBytes)
}
