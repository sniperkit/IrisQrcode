package tools

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
	"os"
)

/**
 * 判断文件是否存在
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 生成32位MD5
 */
func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

/**
 * 生成随机字符串
 */
func GetRandomString(lent int) string{
	str1 := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str2 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes1 := []byte(str1)
	bytes2 := []byte(str2)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result = append(result, bytes1[r.Intn(len(bytes1))])
	for i := 1; i < lent; i++ {
		result = append(result, bytes2[r.Intn(len(bytes2))])
	}
	return string(result)
}
