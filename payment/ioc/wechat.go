package ioc

import (
	"Webook/payment/events"
	"Webook/payment/repository"
	"Webook/payment/service/wechat"
	"Webook/pkg/logger"
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"os"
)

func InitWechatClient(cfg WechatConfig) *core.Client {

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(
		// 注意这个文件我没有上传，所以你需要准备一个
		cfg.KeyPath,
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client
	client, err := core.NewClient(
		ctx,
		option.WithWechatPayAutoAuthCipher(
			cfg.MchID, cfg.MchSerialNum,
			mchPrivateKey, cfg.MchKey),
	)
	if err != nil {
		panic(err)
	}
	return client
}

func InitWechatNativeService(
	cli *core.Client,
	repo repository.PaymentRepository,
	l logger.Logger,
	producer events.Producer,
	cfg WechatConfig) *wechat.NativePaymentService {
	return wechat.NewNativePaymentService(&native.NativeApiService{
		Client: cli,
	}, repo, producer, l, cfg.AppID, cfg.MchID)
}

func InitWechatNotifyHandler(cfg WechatConfig) *notify.Handler {
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(cfg.MchID)
	// 3. 使用apiv3 key、证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(cfg.MchKey,
		verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		panic(err)
	}
	return handler
}

func InitWechatConfig() WechatConfig {
	return WechatConfig{
		AppID:        os.Getenv("WEPAY_APP_ID"),
		MchID:        os.Getenv("WEPAY_MCH_ID"),
		MchKey:       os.Getenv("WEPAY_MCH_KEY"),
		MchSerialNum: os.Getenv("WEPAY_MCH_SERIAL_NUM"),
		CertPath:     "./config/cert/apiclient_cert.pem",
		KeyPath:      "./config/cert/apiclient_key.pem",
	}
}

type WechatConfig struct {
	AppID        string
	MchID        string
	MchKey       string
	MchSerialNum string

	// 证书
	CertPath string
	KeyPath  string
}
