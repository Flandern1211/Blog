package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// 使用 bcrypt 对字符串进行加密生成一个哈希值
func BcryptHash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(bytes), err
}

// 使用 bcrypt 对比 明文字符串 和 哈希值
func BcryptCheck(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

func MD5(str string, b ...byte) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(b))
}

// RandomCode 生成指定长度的随机数字验证码
func RandomCode(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result[i] = digits[num.Int64()]
	}
	return string(result)
}
