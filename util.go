
package glwapi

import (
	"crypto/sha1"
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var rune_len = len(letterRunes)


func R(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	rs := rand.Perm(rune_len)
	for i := range b {
		b[i] = letterRunes[rs[i]]
	}
	return string(b[:n])
}


func Sha1(str string) string {
	ret := sha1.Sum([]byte(str))
	sign := fmt.Sprintf("%x", ret)
	return sign
}


func Md5(str string) string {
	ret := md5.Sum([]byte(str))
	sign := fmt.Sprintf("%x", ret)
	return sign
}
