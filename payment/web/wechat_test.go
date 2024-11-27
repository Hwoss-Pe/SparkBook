package web

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

// 为了验证处理回调的逻辑，最好是手动记下来一个请求
// 这样可以反复运行测试
func TestHandleCallback(t *testing.T) {
	body := `
{"id":"13de4f2b-a78b-5e2f-8773-b30f39b0d341","create_time":"2024-01-11T21:16:33+08:00","resource_type":"encrypt-resource","event_type":"TRANSACTION.SUCCESS","summary":"支付成功","resource":{"original_type":"transaction","algorithm":"AEAD_AES_256_GCM","ciphertext":"UxIDjTS1LqTUqy/6KpsCuqYf+FeBInzOYgjeGsomF6PUmUVjW5XlD0PUSvqZqZtLkbNHJkRLd/HDJYHItpkRy3l8xhx3SGd3B2tvWO0un5w7BBpUo8XetEiXoVagzSpvUDrIKBczckM6ILZRK766dHNx7Z82d10sBAls4eCQwg446gojdU3uOPzM6TxGw8Muf3NJsNh+B0F2tZmwDk0yi4NTgaOwZIPv8Q1l/LmAhbBVLOJCNz28XKGMckQWGyzkTX0X9wiNPTJ3GrEwwPcnwTvERqvFhWsfHK6iewCuAlFs6D4Ta8pi6cs2i2ehBeMgxgNCOYc0LFpIGyd38G3i8+1wCzWFKz3PfswUbNU6w+8elhv8l6YYnnBc/pqq20lwnDpx/DOhTBqMn16TU+cD2rA2UlpMG5NdK+GpysJzZtbCqEg2lgL8xMsHlP4JnIWH6HjVXHMY5tOSN3v6NCuo88/IW6w0UDPbPG4fx5xD8RjTpwEO0nMnHO2JOm2sKfnlQ5pocaGhYUfpvW3AfP54f2bF335l20LctQYnc0B519Lnpqeo6phAgCKdtLl/h2/OuPK/","associated_data":"transaction","nonce":"CkVASomQkPSt"}}
`
	req, err := http.NewRequest(http.MethodPost, "/pay/callback", io.NopCloser(bytes.NewBufferString(body)))
	require.NoError(t, err)
	req.Header.Add("Wechatpay-Serial", "4522CC613021F88AECCFAA187F24BB9A8C73ED61")
	req.Header.Add("Wechatpay-Signature", "pWPMNLXxrIcPTg9R8Tsw9D4RaEb4MBt1eLTaLejXtLctqmAxFYuxABLrR+710HNwMtuupupfok8sixlZMfEcHIdsq14ZTGViEB/5j4bMrBJCAid6Dv/zGuJ6rGlOptAb/+ebLDoshGlI1BL3SnUVBFpOnKmj88Ak4SCnZ9iMq9/HzimxvLkOHWgau+WooKkurRd1r41r1W0wmM0ICgetHwwCMe1jA2D2stCHFZR577+PomDBLT1zUDSAXfKWNetnnukNZqnG4Y4huMd9m6OtfnSfl96UJxhghFn1RDMpBgsDdz6VHuxIAOCmPHYFbCyXiRUYjpP8AhROHc3mxdZEnw==")
	req.Header.Add("Wechatpay-Signature-Type", "WECHATPAY2-SHA256-RSA2048")
	req.Header.Add("Wechatpay-Timestamp", "1704979024")
	req.Header.Add("Wechatpay-Nonce", "uhnrNaJIvLnRlDJaYNGA2Q4rcymFBkXE")
}
