package Libaries

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"meizi/src/Redis"
)

func Md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func CreateRandomString(len int) string  {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < len ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

 func  GetUserIdByToken(token string) []interface{} {
	var ctx = context.Background()
	val, err := Redis.RedisPool.HMGet(ctx, "token_" + token,  "user_id", "nick_name", "username").Result()
	if err != nil {
		panic(err)
	}

	return val
}
