package util

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

const (
	KEY                    string = "JWT_ARY-START"
	DEFAULT_EXPIRE_SECONDS int    = 7200 //默认过期时间
)

//管理员 token数据
type AdminJwtData struct {
	Id       int
	RealName string
	Sex      int8
}

//JWT -- json web token
//HEADER PAYLOAD SIGNATURE
//This struct is the PAYLOAD
type Claims struct {
	Admin *AdminJwtData
	jwt.StandardClaims
}

//刷新jwt token
func RefreshToken(tokenString string, expiredSeconds int) (string, error) {
	if expiredSeconds == 0 {
		expiredSeconds = DEFAULT_EXPIRE_SECONDS
	}
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
			return []byte(KEY), nil
		})
	if token == nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", err
	}
	mySigningKey := []byte(KEY)
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	newClaims := Claims{
		claims.Admin,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    claims.Admin.RealName,
			IssuedAt:  time.Now().Unix(),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodES256, newClaims)
	tokenStr, err := newToken.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("generate new fresh json web token failed !! error :", err)
		return "", err
	}
	return tokenStr, err
}

//验证jwt token
func ValidateToken(tokenString string) (info *AdminJwtData, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return []byte(KEY), nil
		})
	if token == nil {
		return
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		info = claims.Admin
	} else {
		fmt.Println("validate tokenString failed !!!", err)
	}
	return
}

//获取jwt token
func GenerateToken(info *AdminJwtData, expiredSeconds int) (tokenString string, err error) {
	if expiredSeconds == 0 {
		expiredSeconds = DEFAULT_EXPIRE_SECONDS
	}
	mySigningKey := []byte(KEY)
	expireAt := time.Now().Add(time.Second * time.Duration(expiredSeconds)).Unix()
	claims := Claims{
		info,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    info.RealName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("generate json web token failed !! error :", err)
	} else {
		tokenString = tokenStr
	}
	return
}
