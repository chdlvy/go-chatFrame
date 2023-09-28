package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 加密解密的秘钥
var mySigningKey []byte = []byte("asdfghjkl$")

// 加密
func GetToken(content map[string]interface{}) string {
	// map中的内容转到jwt.MapClaims中
	mapClaims := make(jwt.MapClaims)
	for k, v := range content {
		mapClaims[k] = v
	}
	// 设置默认时间
	_, ok := mapClaims["exp"]
	if !ok {
		mapClaims["exp"] = time.Now().Unix() + 50
	}
	//new一个token，参数1：加密算法，参数2：需要加密的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	// 加密
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("加密失败：", err.Error())
		return "加密失败：" + err.Error()
	}
	return tokenString
}

// 解密
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	// fmt.Println("解密：", token)
	// fmt.Println("token.Claims：", token.Claims)
	// fmt.Println("name：", token.Claims.(jwt.MapClaims)["name"])
	return token.Claims.(jwt.MapClaims), nil
}
