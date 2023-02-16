package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (pc *CloudrevePayController) BearerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			logrus.WithField("Authorization", authorization).Debugln("Authorization 头缺失或无效")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"data":    "",
				"message": "Authorization 头缺失或无效",
			})
			return
		}

		authorizations := strings.Split(strings.TrimPrefix(authorization, "Bearer "), ":")
		if len(authorizations) != 2 {
			logrus.WithField("Authorization", authorization).WithField("len.auth", len(authorizations)).Debugln("Authorization 头无效")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"data":    "",
				"message": "Authorization 头无效",
			})
			return
		}

		// 验证是否过期
		signature := strings.TrimPrefix(authorization, "Bearer ")
		expires, err := strconv.ParseInt(authorizations[1], 10, 64)
		if err != nil {
			logrus.WithField("Authorization", authorization).WithField("ttlUnix", authorizations[1]).Debugln("Authorization 头无效，无法解析 ttl")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"data":    "",
				"message": "Authorization 头无效，无法解析 ttl",
			})
			return
		}

		// 如果签名过期
		if expires < time.Now().Unix() && expires != 0 {
			logrus.WithField("Authorization", authorization).WithField("ttlUnix", authorizations[1]).Debugln("Authorization 头无效，签名已过期")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"data":    "",
				"message": "Authorization 头无效，签名已过期",
			})
			return
		}

		auth := &HMACAuth{
			CloudreveKey: []byte(pc.Conf.CloudreveKey),
		}

		if generatedSign := auth.Sign(getSignContent(c.Request), expires); signature != generatedSign {
			logrus.WithField("Authorization", authorization).WithField("generatedSign", generatedSign).Debugln("Authorization 头无效，签名不匹配")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"data":    "",
				"message": "Authorization 头无效，签名不匹配",
			})
			return
		}
	}
}
