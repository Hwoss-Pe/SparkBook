package wechat

import (
	"Webook/payment/domain"
	"Webook/payment/ioc"
	"Webook/payment/repository"
	"Webook/payment/repository/dao"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ecodeclub/ekit/net/httpx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

// 下单
func TestNativeService_Prepay(t *testing.T) {
	appid := os.Getenv("WEPAY_APP_ID")
	mchID := os.Getenv("WEPAY_MCH_ID")
	mchKey := os.Getenv("WEPAY_MCH_KEY")
	mchSerialNumber := os.Getenv("WEPAY_MCH_SERIAL_NUM")
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(
		"/Users/webook/payment/config/cert/apiclient_key.pem",
	)
	require.NoError(t, err)
	ctx := context.Background()
	// 使用商户私钥等初始化 client
	client, err := core.NewClient(
		ctx,
		option.WithWechatPayAutoAuthCipher(mchID, mchSerialNumber, mchPrivateKey, mchKey),
	)
	require.NoError(t, err)
	nativeSvc := &native.NativeApiService{
		Client: client,
	}
	logger := ioc.InitLogger()
	db := ioc.InitDB()
	repo := repository.NewPaymentRepository(dao.NewPaymentGORMDAO(db))
	kafka := ioc.InitKafka()
	producer := ioc.InitProducer(kafka)
	svc := NewNativePaymentService(nativeSvc, repo, producer, logger, appid, mchID)
	codeUrl, err := svc.Prepay(ctx, domain.Payment{
		Amt: domain.Amount{
			Currency: "CNY",
			Total:    1,
		},
		BizTradeNO:  "test",
		Description: "面试官AI",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, codeUrl)
	t.Log(codeUrl)
}

func TestServer(t *testing.T) {
	http.HandleFunc("/", func(
		writer http.ResponseWriter,
		request *http.Request) {
		writer.Write([]byte("hello, 我进来了"))
	})
	http.ListenAndServe(":8080", nil)
}

func TestGetOpenID(t *testing.T) {
	// 0a3eg4000TDOmR1s4i000ZXLmo1eg40N
	//0f3y4hGa13vxGG0c0pJa1JPVLh4y4hGD
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp := httpx.NewRequest(ctx, http.MethodGet,
		"https://api.weixin.qq.com/sns/jscode2session").
		AddParam("grant_type", "authorization_code").
		AddParam("js_code", "0f3y4hGa13vxGG0c0pJa1JPVLh4y4hGD").
		AddParam("secret", "ac35caf51b0948f472e19b27505deab8").
		AddParam("appid", "wxc6daad0e647d1544").Do()
	fmt.Println(resp)
}

// 支付结果回调的处理
func PayResHandler(w http.ResponseWriter, r *http.Request) {
	// 验签
	var (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775************************************" // 商户证书序列号
		mchAPIV3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)
	ctx := context.Background()
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		log.Fatal("load merchant private key error")
	}
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIV3Key)
	if err != nil {
		log.Println(err)
		return
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(mchAPIV3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), r, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		log.Println(err)
		return
	}
	//处理数据
	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)
	// 返回数据
	val, _ := json.Marshal(Res{
		Code:    "SUCCESS",
		Message: "",
	})
	_, err = w.Write(val)
	if err != nil {
		log.Println(err)
	}
}

type Res struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 根据订单号查询订单
func ExampleJsapiApiService_QueryOrderById() {
	var (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775************************************" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := jsapi.JsapiApiService{Client: client}
	resp, result, err := svc.QueryOrderById(ctx,
		jsapi.QueryOrderByIdRequest{
			// 订单号
			TransactionId: core.String("TransactionId_example"),
			Mchid:         core.String("Mchid_example"),
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call QueryOrderById err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
	}
}

// 关闭订单
func ExampleJsapiApiService_CloseOrder() {
	var (
		mchID                      string = "190000****"                               // 商户号
		mchCertificateSerialNumber string = "3775************************************" // 商户证书序列号
		mchAPIv3Key                string = "2ab9****************************"         // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath("/path/to/merchant/apiclient_key.pem")
	if err != nil {
		log.Print("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
	}

	svc := jsapi.JsapiApiService{Client: client}
	result, err := svc.CloseOrder(ctx,
		jsapi.CloseOrderRequest{
			OutTradeNo: core.String("OutTradeNo_example"),
			Mchid:      core.String("1230000109"),
		},
	)
	if err != nil {
		// 处理错误
		log.Printf("call CloseOrder err:%s", err)
	} else {
		// 处理返回结果
		log.Printf("status=%d", result.Response.StatusCode)
	}
}
