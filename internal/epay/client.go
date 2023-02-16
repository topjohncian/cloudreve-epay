package epay

type Config struct {
	PartnerID string
	Key       string
	Endpoint  string
}

var _ Client = &EPayClient{}

// Client 易支付API
type Client interface {
	// Purchase 生成支付链接和参数
	Purchase(args *PurchaseArgs) (string, map[string]string)
	// Verify 验证回调参数是否符合签名
	Verify(params map[string]string) (*VerifyRes, error)
}

type EPayClient struct {
	Config *Config
}

func NewClient(config *Config) *EPayClient {
	return &EPayClient{
		Config: config,
	}
}
