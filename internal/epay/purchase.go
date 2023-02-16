package epay

import (
	"net/url"

	"github.com/mitchellh/mapstructure"
)

type PurchaseType string

var (
	Alipay PurchaseType = "alipay"
	Wxpay  PurchaseType = "wxpay"
)

type DeviceType string

var (
	PC     DeviceType = "pc"
	MOBILE DeviceType = "mobile"
)

type PurchaseArgs struct {
	// 支付类型
	Type PurchaseType
	// 商家订单号
	ServiceTradeNo string
	// 商品名称
	Name string
	// 金额
	Money string
	// 设备类型
	Device    DeviceType
	NotifyUrl *url.URL
	ReturnUrl *url.URL
}

// const PURCHASE_API = "https://payment.moe/submit.php"

// Purchase 生成支付链接和参数
func (c *EPayClient) Purchase(args *PurchaseArgs) (string, map[string]string) {
	// see https://payment.moe/doc.html
	requestParams := map[string]string{
		"pid":          c.Config.PartnerID,
		"type":         string(args.Type),
		"out_trade_no": args.ServiceTradeNo,
		"notify_url":   args.NotifyUrl.String(),
		"name":         args.Name,
		"money":        args.Money,
		"device":       string(args.Device),
		"sign_type":    "MD5",
		"return_url":   args.ReturnUrl.String(),
		"sign":         "",
	}

	return c.Config.Endpoint, GenerateParams(requestParams, c.Config.Key)
}

const TRADE_SUCCESS = "TRADE_SUCCESS"

type VerifyRes struct {
	// 支付类型
	Type PurchaseType
	// 易支付订单号
	TradeNo string `mapstructure:"trade_no"`
	// 商家订单号
	ServiceTradeNo string `mapstructure:"out_trade_no"`
	// 商品名称
	Name string
	// 金额
	Money string
	// 订单支付状态
	TradeStatus string `mapstructure:"trade_status"`
	// 签名检验
	VerifyStatus bool `mapstructure:"-"`
}

func (c *EPayClient) Verify(params map[string]string) (*VerifyRes, error) {
	sign := params["sign"]
	var verifyRes VerifyRes
	// 从 map 映射到 struct 上
	err := mapstructure.Decode(params, &verifyRes)
	// 验证签名
	verifyRes.VerifyStatus = sign == GenerateParams(params, c.Config.Key)["sign"]
	if err != nil {
		return nil, err
	} else {
		return &verifyRes, nil
	}
}
