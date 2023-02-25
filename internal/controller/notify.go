package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/topjohncian/cloudreve-pro-epay/internal/epay"
)

type NotifyResponse struct {
	Code  int    `json:"code"`
	Error string `json:"eror"`
}

func (pc *CloudrevePayController) Notify(c *gin.Context) {
	query := c.Request.URL.Query()
	params := lo.Reduce(lo.Keys(query), func(r map[string]string, t string, i int) map[string]string {
		r[t] = query.Get(t)
		return r
	}, map[string]string{})

	if epay.GenerateSign(params, pc.Conf.EpayKey) != params["sign"] {
		c.String(400, "fail")
		logrus.Warningln("签名验证失败")
		return
	}

	orderId := c.Param("id")
	if orderId == "" {
		logrus.Debugln("无效的订单号")
		c.String(400, "fail")
		return
	}

	request, ok := pc.Cache.Get(PurchaseSessionPrefix + orderId)
	if !ok {
		logrus.WithField("id", orderId).Debugln("订单信息不存在")
		c.String(400, "fail")
		return
	}

	order, ok := request.(*PurchaseRequest)
	if !ok {
		logrus.WithField("id", orderId).Debugln("订单信息非法")
		c.String(400, "fail")
		return
	}

	if params["trade_status"] == "TRADE_SUCCESS" {
		amount := decimal.NewFromInt(int64(order.Amount)).Div(decimal.NewFromInt(100))
		realAmount, err := decimal.NewFromString(params["money"])
		if err != nil {
			logrus.WithError(err).WithField("id", orderId).Debugln("无法解析订单金额")
			c.String(400, "fail")
			return
		}
		if !realAmount.Equal(amount) {
			logrus.WithField("id", orderId).Debugln("订单金额不符")
			c.String(400, "fail")
			return
		}

		err = retry.Do(func() error {
			var notifyRes NotifyResponse
			resp, err := pc.Client.R().SetSuccessResult(&notifyRes).Get(order.NotifyUrl)
			if err != nil {
				logrus.WithField("id", orderId).WithError(err).Errorln("通知失败")
				return err
			}

			if !resp.IsSuccessState() {
				logrus.WithField("id", orderId).WithField("dump", resp.Dump()).Errorln("通知失败")
				return errors.New("http code: " + strconv.Itoa(resp.StatusCode))
			}

			if notifyRes.Code != 0 {
				logrus.WithField("id", orderId).WithField("dump", resp.Dump()).WithField("error", notifyRes.Error).Errorln("通知失败")
				return errors.New("code: " + strconv.Itoa(notifyRes.Code) + ", error: " + notifyRes.Error)
			}

			return nil
		}, retry.Attempts(5), retry.Delay(10), retry.OnRetry(func(n uint, err error) {
			logrus.WithField("id", orderId).WithField("n", n).WithError(err).Infoln("通知失败，重试")
		}))

		if err != nil {
			logrus.WithField("id", orderId).WithError(err).Errorln("通知失败")
			c.String(400, "fail")
			return
		}

		logrus.WithField("id", orderId).Infoln("通知成功")
		c.String(200, "success")

		pc.Cache.Delete([]string{orderId}, PurchaseSessionPrefix)
		return
	}

	c.String(200, "success")
}

func (pc *CloudrevePayController) Return(c *gin.Context) {
	c.HTML(http.StatusOK, "return.tmpl", gin.H{})
}
