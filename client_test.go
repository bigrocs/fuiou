package fuiou

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bigrocs/fuiou/requests"
)

func TestPay(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.pay"
	request.BizContent = map[string]interface{}{
		"order_type":             "ALIPAY",                // 商户、部门编号
		"goods_des":              "测试商品",                  // 被扫付款码
		"mchnt_order_no":         "513457061273811818896", // 商户订单号(商户交易系统中唯一)
		"order_amt":              "1",
		"term_ip":                "127.0.0.1",
		"txn_begin_ts":           time.Now().Format("20060102150405"),
		"auth_code":              "289200593914815731",
		"reserved_terminal_info": `{"serial_num":"12345678901SN012"}`,
		"goods_detail":           "",
		"addn_inf":               "",
		"curr_type":              "",
		"goods_tag":              "",
		"sence":                  "",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestPay", r, err)
	t.Log(r, err, "|||")
}

func TestPayQuery(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.orderquery"
	request.BizContent = map[string]interface{}{
		"order_type":     "ALIPAY",                // 商户、部门编号
		"mchnt_order_no": "513457061273811818895", // 商户订单号(商户交易系统中唯一)
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestPayQuery", r, err)
	t.Log(r, err, "|||")
}

func TestPayRefund(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.refund"
	request.BizContent = map[string]interface{}{
		"order_type":      "ALIPAY",                // 商户、部门编号
		"mchnt_order_no":  "513457061273811818896", // 商户订单号(商户交易系统中唯一)
		"refund_order_no": "513457061273811818896-1",
		"total_amt":       "1",
		"refund_amt":      "1",
		"operator_id":     "",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestPayRefund", r, err)
	t.Log(r, err, "|||")
}

func TestPayRefundQuery(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.refundquery"
	request.BizContent = map[string]interface{}{
		"refund_order_no": "513457061273811818896-1",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestPayRefundQuery", r, err)
	t.Log(r, err, "|||")
}

func TestPayQrcode(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.qrcode"
	request.BizContent = map[string]interface{}{
		"order_type":     "ALIPAY",                // 商户、部门编号
		"goods_des":      "测试商品",                  // 被扫付款码
		"mchnt_order_no": "513457061273811818898", // 商户订单号(商户交易系统中唯一)
		"order_amt":      "1",
		"term_ip":        "127.0.0.1",
		"txn_begin_ts":   time.Now().Format("20060102150405"),
		"goods_detail":   "",
		"addn_inf":       "",
		"curr_type":      "",
		"goods_tag":      "",
		"notify_url":     "http://127.0.0.1",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestPayQrcode", r, err)
	t.Log(r, err, "|||")
}

func TestQueryOpenid(t *testing.T) {
	// 创建连接
	// client := NewClient()
	// client.Config.AppId = "10000000000003188540"
	// client.Config.FuiouPublicKey = os.Getenv("PAY_Fuiou_Fuiou_PUBLIC_KEY")
	// client.Config.PrivateKey = os.Getenv("PAY_Fuiou_PRIVATE_KEY")
	// client.Config.ReturnSignType = "RSA"
	// // client.Config.Sandbox = true
	// // 配置参数
	// request := requests.NewCommonRequest()
	// request.ApiName = "pay.openid"
	// request.BizContent = map[string]interface{}{
	// 	"mersepNo":  "161504043788", // 商户、部门编号
	// 	"oauthCode": "284959078663380689",
	// 	// "subAppId":  "wxb6f76f399a3e7561",
	// }
	// // 请求
	// response, err := client.ProcessCommonRequest(request)
	// if err != nil {
	// 	fmt.Println(response, err)
	// }
	// r, err := response.GetVerifySignDataMap()
	// fmt.Println("TestQueryOpenid", r, err)
	// t.Log(r, err, "|||")
}

func TestJsApi(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.jsapi"
	request.BizContent = map[string]interface{}{
		"goods_des":      "测试商品", // 被扫付款码
		"goods_detail":   "",
		"goods_tag":      "",
		"product_id":     "",
		"addn_inf":       "",
		"mchnt_order_no": "513457061273811818910", // 商户订单号(商户交易系统中唯一)
		"curr_type":      "",
		"order_amt":      "1",
		"term_ip":        "127.0.0.1",
		"txn_begin_ts":   time.Now().Format("20060102150405"),
		"notify_url":     "http://127.0.0.1",
		"limit_pay":      "",
		"trade_type":     "JSAPI", // 商户、部门编号
		"openid":         "",
		"sub_openid":     "1212112",
		"sub_appid":      "12211212",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestJsApi", r, err)
	t.Log(r, err, "|||")
}
func TestOpenid(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.InsCd = os.Getenv("PAY_FUIOU_INS_CD")
	client.Config.FuiouPublicKey = os.Getenv("PAY_FUIOU_FUIOU_PUBLIC_KEY")
	client.Config.PrivateKey = os.Getenv("PAY_FUIOU_PRIVATE_KEY")
	client.Config.MchntCd = os.Getenv("PAY_FUIOU_MCHNT_CD")
	// client.Config.Sandbox = true
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "pay.openid"
	request.BizContent = map[string]interface{}{
		"auth_code":  "133149051993393767",
		"term_ip":    "127.0.0.1",
		"sub_appid":  "",
		"order_type": "WECHAT",
		"order_amt":  "",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestOpenid", r, err)
	t.Log(r, err, "|||")
}
