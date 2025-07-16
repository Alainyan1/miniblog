// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package token

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

// token的配置选项
type Config struct {
	// 用于签发和解析token的密钥
	key string
	// token中用户的key
	identityKey string
	// 签发的token过期时间
	expiration time.Duration
}

var (
	// 默认值
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey", 2 * time.Hour}
	once   sync.Once
)

// 设置包级别的配置config, config会用于本包后面的token签发和解析
func Init(key string, identityKey string, expiration time.Duration) {
	once.Do(func() {
		if key != "" {
			config.key = key // 设置密钥
		}
		if identityKey != "" {
			config.identityKey = identityKey // 设置身份键
		}
		if expiration != 0 {
			config.expiration = expiration
		}
	})
}

// Parse函数用于解析JWT字符串并提取用户身份
// 使用指定密钥key解析token, 解析成功返回token上下文, 否则报错
func Parse(tokenString string, key string) (string, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保token加密算法是预期加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	var identityKey string
	// 如果解析成功, 从token中取出token的主题, 这里的claims是The second segment of the token, 即payload
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if key, exists := claims[config.identityKey]; exists {
			if identity, valid := key.(string); valid {
				identityKey = identity // 获取身份键
			}
		}
	}

	if identityKey == "" {
		return "", jwt.ErrSignatureInvalid
	}

	return identityKey, nil
}

// 从请求中获取JWT, 将其传递给Parse函数来解析
func ParseRequest(ctx context.Context) (string, error) {
	var (
		token string
		err   error
	)

	// 提取ctx的实际类型, 仅用于switch语句
	switch typed := ctx.(type) {
	// gin开发的http服务
	case *gin.Context:
		header := typed.Request.Header.Get("Authorization")
		if len(header) == 0 {
			//nolint: err113
			return "", errors.New("the length of the `Authorization` header is zero")
		}
		// fmt.Sscanf 用于从字符串中解析格式化数据
		_, _ = fmt.Sscanf(header, "Bearer %s", &token) // 解析Bearer token
	// grpc服务
	default:
		token, err = auth.AuthFromMD(typed, "Bearer")
		if err != nil {
			return "", status.Errorf(codes.Unauthenticated, "invalid auth token")
		}
	}

	return Parse(token, config.key) // 解析token
}

// Sign 使用 jwtSecret 签发 token, token 的 claims 中会存放传入的 subject.
func Sign(identityKey string) (string, time.Time, error) {
	// 计算过期时间
	expireAt := time.Now().Add(config.expiration)

	// token内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,       // 用户身份
		"nbf":              time.Now().Unix(), // token生效时间
		"iat":              time.Now().Unix(), // 签发时间
		"exp":              expireAt.Unix(),   // 过期时间
	})

	// 签发token
	tokenString, err := token.SignedString([]byte(config.key))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expireAt, nil
}
