#!/bin/bash

# Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/Alainyan1/miniblog.

# 生成jwt token

# 定义header, header包含两部分信息: token的类型typ和使用的算法alg
HEADER='{"alg":"H256","typ":"JWT"}'

# 定义Payload
# payload标准字段: iss: token的签发者, sub: 主题, exp: token过期时间, aud: 接受token的一方, iat: token的签发时间, nbf: token的生效时间, jti: token的唯一标识
PAYLOAD='{"exp":1739078005,"iat":1735478005,"nbf":1735478005,"x-user-id":"user-w6irkg"}'

# 定义Secret(用于签名), Sigature的生成方式如下: 使用Base64对header.payload进行编码, 使用Secret对编码后的内容进行加密, 加密后的内容即为Signature
SECRET="Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5"

# 1. Base64编码Header
# -n避免添加换行符
# tr -d '='删除所有 = 字符, 标准Base64编码中, = 用于填充使输出长度为4的倍数. 移除 = 得到无填充的Base64变体, 常用于URL安全的场景
# tr命令进行字符串替换, 将所有 / 替换为 _ 将所有 + 替换为 -, / 和 + 在URL中有特殊含义, 而 _ 和 - 更适合URL使用
# tr -d '\n'移除换行符
HEADER_BASE64=$(echo -n "${HEADER}" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 2. Base64编码Payload
PAYLOAD_BASE64=$(echo -n "${PAYLOAD}" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 3. 拼接Header和Payload为签名数据
SIGNING_INPUT="${HEADER_BASE64}.${PAYLOAD_BASE64}"

# 4. 使用HMAC SHA256算法生成签名
SIGNATURE=$(echo -n "${SIGNING_INPUT}" | openssl dgst -sha256 -hmac "${SECRET}" -binary | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 5. 拼接最终的JWT Token
JWT="${SIGNING_INPUT}.${SIGNATURE}"

# 输出JWT Token
echo "Generated JWT Token:"
echo "${JWT}"