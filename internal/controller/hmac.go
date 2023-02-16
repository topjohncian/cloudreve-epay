package controller

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const CrHeaderPrefix = "X-Cr-"

// General 通用的认证接口
var General Auth

// Auth 鉴权认证
type Auth interface {
	// 对给定Body进行签名,expires为0表示永不过期
	Sign(body string, expires int64) string
	// 对给定Body和Sign进行检查
	Check(body string, sign string) error
}

// SignRequest 对PUT\POST等复杂HTTP请求签名，只会对URI部分、
// 请求正文、`X-Cr-`开头的header进行签名
func SignRequest(instance Auth, r *http.Request, expires int64) *http.Request {
	// 处理有效期
	if expires > 0 {
		expires += time.Now().Unix()
	}

	// 生成签名
	sign := instance.Sign(getSignContent(r), expires)

	// 将签名加到请求Header中
	r.Header["Authorization"] = []string{"Bearer " + sign}
	return r
}

// // CheckRequest 对复杂请求进行签名验证
// func CheckRequest(instance Auth, r *http.Request) error {
// 	var (
// 		sign []string
// 		ok   bool
// 	)
// 	if sign, ok = r.Header["Authorization"]; !ok || len(sign) == 0 {
// 		return ErrAuthHeaderMissing
// 	}
// 	sign[0] = strings.TrimPrefix(sign[0], "Bearer ")

// 	return instance.Check(getSignContent(r), sign[0])
// }

// RequestRawSign 待签名的HTTP请求
type RequestRawSign struct {
	Path   string
	Header string
	Body   string
}

// NewRequestSignString 返回JSON格式的待签名字符串
func NewRequestSignString(path, header, body string) string {
	req := RequestRawSign{
		Path:   path,
		Header: header,
		Body:   body,
	}
	res, _ := json.Marshal(req)
	return string(res)
}

// getSignContent 签名请求 path、正文、以`X-`开头的 Header. 如果请求 path 为从机上传 API，
// 则不对正文签名。返回待签名/验证的字符串
func getSignContent(r *http.Request) (rawSignString string) {
	// 读取所有body正文
	var body = []byte{}
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	// 决定要签名的header
	var signedHeader []string
	for k := range r.Header {
		if strings.HasPrefix(k, CrHeaderPrefix) && k != CrHeaderPrefix+"Filename" {
			signedHeader = append(signedHeader, fmt.Sprintf("%s=%s", k, r.Header.Get(k)))
		}
	}
	sort.Strings(signedHeader)

	// 读取所有待签名Header
	rawSignString = NewRequestSignString(r.URL.Path, strings.Join(signedHeader, "&"), string(body))

	return rawSignString
}

// HMACAuth HMAC算法鉴权
type HMACAuth struct {
	CloudreveKey []byte
}

// Sign 对给定Body生成expires后失效的签名，expires为过期时间戳，
// 填写为0表示不限制有效期
func (auth HMACAuth) Sign(body string, expires int64) string {
	h := hmac.New(sha256.New, auth.CloudreveKey)
	expireTimeStamp := strconv.FormatInt(expires, 10)
	_, err := io.WriteString(h, body+":"+expireTimeStamp)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)) + ":" + expireTimeStamp
}

// // Check 对给定Body和Sign进行鉴权，包括对expires的检查
// func (auth HMACAuth) Check(body string, sign string) error {
// 	signSlice := strings.Split(sign, ":")
// 	// 如果未携带expires字段
// 	if signSlice[len(signSlice)-1] == "" {
// 		return ErrExpiresMissing
// 	}

// 	// 验证是否过期
// 	expires, err := strconv.ParseInt(signSlice[len(signSlice)-1], 10, 64)
// 	if err != nil {
// 		return ErrAuthFailed.WithError(err)
// 	}
// 	// 如果签名过期
// 	if expires < time.Now().Unix() && expires != 0 {
// 		return ErrExpired
// 	}

// 	// 验证签名
// 	if auth.Sign(body, expires) != sign {
// 		return ErrAuthFailed
// 	}
// 	return nil
// }
