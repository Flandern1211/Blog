package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

// Format 格式化字符串（去空格）
func Format(s string) string {
	return strings.TrimSpace(s)
}

// GenEmailVerificationInfo 生成邮箱验证信息（用户名:密码 的 base64）
func GenEmailVerificationInfo(username, password string) string {
	data := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// ParseEmailVerificationInfo 解析邮箱验证信息
func ParseEmailVerificationInfo(info string) (string, string, error) {
	data, err := base64.StdEncoding.DecodeString(info)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid verification info")
	}
	return parts[0], parts[1], nil
}
