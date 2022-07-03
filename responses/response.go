/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless resuired by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package responses

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/micro/go-micro/v2/util/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/bigrocs/fuiou/config"
	"github.com/bigrocs/fuiou/requests"
	"github.com/bigrocs/fuiou/util"
)

const (
	CLOSED     = "CLOSED"     // -1 订单关闭
	USERPAYING = "USERPAYING" // 0	订单支付中
	SUCCESS    = "SUCCESS"    // 1	订单支付成功
	WAITING    = "WAITING"    // 2	系统执行中请等待
)

// CommonResponse 公共回应
type CommonResponse struct {
	Config      *config.Config
	Request     *requests.CommonRequest
	httpContent []byte
	json        string
}

type Map *mxj.Map

// NewCommonResponse 创建新的请求返回
func NewCommonResponse(config *config.Config, request *requests.CommonRequest) (response *CommonResponse) {
	c := &CommonResponse{}
	c.Config = config
	c.Request = request
	return c
}

// GetHttpContentJson 获取 JSON 数据
func (res *CommonResponse) GetHttpContentJson() string {
	return res.json
}

// GetHttpContentMap 获取 MAP 数据
func (res *CommonResponse) GetHttpContentMap() (mxj.Map, error) {
	return mxj.NewMapJson([]byte(res.json))
}

// GetVerifySignDataMap 获取 GetVerifySignDataMap 校验后数据数据
func (res *CommonResponse) GetVerifySignDataMap() (m mxj.Map, err error) {
	r, err := res.GetHttpContentMap()
	if err != nil {
		return r, err
	}
	if r["sign"] != nil {
		ok, err := util.VerifySign(res.GetSignData(), r["sign"].(string), res.Config.FuiouPublicKey, "MD5")
		if err != nil {
			return r, err
		}
		if ok {
			return res.GetSignDataMap()
		}
	} else {
		return r, errors.New("res sign is not")
	}
	return
}

// GetSignData 获取 SignData 数据
func (res *CommonResponse) GetSignData() string {
	return res.json
}

// GetSign 获取 Sign 数据
func (res *CommonResponse) GetSign() (string, error) {
	mv, err := res.GetHttpContentMap()
	if err != nil {
		return "", err
	}
	if _, ok := mv["sign"]; ok { //去掉 xml 外层
		return mv["sign"].(string), err
	}
	return "", err
}

// SetHttpContent 设置请求信息
func (res *CommonResponse) SetHttpContent(httpContent []byte, dataType string) {
	res.httpContent = httpContent
	switch dataType {
	case "xml":
		con, err := url.QueryUnescape(string(res.httpContent))
		if err != nil {
			log.Fatal(err)
		}
		r, err := res.GbkToUtf8([]byte(con))
		if err != nil {
			log.Fatal(err)
		}
		mv, err := mxj.NewMapXml(r) // unmarshal
		if err != nil {
			log.Fatal(err)
		}
		var str interface{}
		if _, ok := mv["xml"]; ok { //去掉 xml 外层
			str = mv["xml"]
		} else {
			str = mv
		}
		jsonStr, _ := json.Marshal(str)
		res.json = string(jsonStr)
		log.Info("Fuiou[PostForm]res", res.json)
	case "string":
		res.json = string(res.httpContent)
	}
}
func (res *CommonResponse) GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	str := string(d)
	str = strings.ReplaceAll(str, `<?xml version="1.0" encoding="gbk" standalone="yes"?>`, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	str = strings.ReplaceAll(str, `<?xml version="1.0" encoding="GBK" standalone="yes"?>`, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	str = strings.ReplaceAll(str, `<?xml version="1.0" encoding="GB2312" standalone="yes" ?>`, `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>`)
	return []byte(str), nil
}

// GetSignDataMap 获取 MAP 数据
func (res *CommonResponse) GetSignDataMap() (mxj.Map, error) {
	data := mxj.New()
	content, err := mxj.NewMapJson([]byte(res.GetSignData()))
	if err != nil {
		return nil, err
	}
	if res.Request.ApiName == "pay.pay" {
		data = res.handerFuiouTradePay(content)
	}
	if res.Request.ApiName == "pay.orderquery" {
		data = res.handerFuiouTradeQuery(content)
	}
	if res.Request.ApiName == "pay.refund" {
		data = res.handerFuiouTradeRefund(content)
	}
	if res.Request.ApiName == "pay.refundquery" {
		data = res.handerFuiouTradeRefundQuery(content)
	}
	if res.Request.ApiName == "pay.qrcode" {
		data = res.handerFuiouTradeQrcode(content)
	}
	if res.Request.ApiName == "pay.openid" {
		data = res.handerFuiouTradeOpenid(content)
	}
	if res.Request.ApiName == "pay.jsapi" {
		data = res.handerFuiouTradeJsApi(content)
	}
	if res.Request.ApiName == "pay.openid" {
		data = res.handerFuiouTradeOpenid(content)
	}
	data["channel"] = "fuiou" //渠道
	data["content"] = content
	return data, err
}

// alipay_open_id:2088002104076813
// bank_trade_no:
//  channel:fuiou
//  content:map[addn_inf:
//  buyer_id:2088002104076813
//  ins_cd:08A9999999
//  mchnt_cd:0002900F0370542
//   order_type:ALIPAY
//   random_str:W8OBJ3VK928C9QYFQDZDP65UJERXAGCR
//   reserved_bank_type:BANKCARD
//   reserved_buyer_logon_id:big***@qq.com
//   reserved_channel_order_id:513457061273811818895
//   reserved_coupon_fee:
//   reserved_fund_bill_list:[{"amount":"0.01","fund_channel":"BANKCARD",
//   "fund_type":"CREDIT_CARD"}]
//   reserved_fy_order_no:
//   reserved_fy_settle_dt:20220701
//   reserved_fy_trace_no:070005041733
//   reserved_is_credit:1
//   reserved_mchnt_order_no:513457061273811818895
//   reserved_risk_info:
//   reserved_settlement_amt:1
//   reserved_txn_fin_ts:20220701180849
//   result_code:000000
//   result_msg:SUCCESS sign:
//   +
//   T4= term_id:
//   total_amount:1
//   transaction_id:2022070122001476811446553236]
//   out_trade_no:513457061273811818895
//   return_code:SUCCESS
//   return_msg:SUCCESS
//   status:SUCCESS time_end:513457061273811818895
//    total_fee:1 trade_no:2022070122001476811446553236]

// handerFuiouTradePay
func (res *CommonResponse) handerFuiouTradePay(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["status"] = "" // 状态
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		data["return_code"] = SUCCESS
		data["status"] = SUCCESS
		switch content["result_msg"] {
		case "USERPAYING":
			data["status"] = USERPAYING
		case "CLOSED":
			data["status"] = CLOSED
		case "REVOKED":
			data["status"] = CLOSED
		case "NOTPAY":
			data["status"] = USERPAYING
		}
		total_amt, _ := strconv.ParseInt(content["total_amount"].(string), 10, 64)
		data["total_fee"] = total_amt
		if v, ok := content["reserved_fy_order_no"]; ok {
			data["bank_trade_no"] = v // 富有订单
		}
		data["trade_no"] = content["transaction_id"]
		if v, ok := content["reserved_mchnt_order_no"]; ok {
			data["out_trade_no"] = v
		}
		if v, ok := content["reserved_txn_fin_ts"]; ok && v != "" {
			data["time_end"] = v
		}
		// channel
		if v, ok := content["order_type"]; ok {
			switch v {
			case "WECHAT":
				data["wechat_open_id"] = content["buyer_id"]
			case "ALIPAY":
				data["alipay_open_id"] = content["buyer_id"]
			}
		}

	} else {
		data["return_code"] = "FAIL"
		data["status"] = CLOSED
		if content["result_code"] == "030010" {
			data["status"] = USERPAYING
		}
		if content["result_code"] == "010002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "9999" {
			data["status"] = WAITING
		}
		if content["result_code"] == "010001" {
			data["status"] = WAITING
		}
		if content["result_code"] == "2002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030003" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030004" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1013" {
			data["status"] = WAITING
		}
	}
	return data
}

// addn_inf:
// buyer_id:2088002104076813
// ins_cd:08A9999999
// mchnt_cd:0002900F0370542
// mchnt_order_no:513457061273811818895
// order_amt:1
// order_type:ALIPAY
// random_str:W73KI8UQUND6MGVCEOSF7DD49YO1HLF9
// reserved_bank_type:BANKCARD
// reserved_buyer_logon_id:big***@qq.com
// reserved_channel_order_id:513457061273811818895
// reserved_coupon_fee:
// reserved_fund_bill_list:
// reserved_fy_settle_dt:20220701
// reserved_fy_trace_no:070005041733
// reserved_is_credit:1
// reserved_risk_info:
// reserved_settlement_amt:1
// reserved_txn_fin_ts:20220701180849
// result_code:000000 result_msg:SUCCESS
// sign:YLB+gNRm4ebslB47GFPCAu4YhrvbljlOamAAjMuTFqv/xR3mVMQHbElT4+e2473NelNHpJguaGFRLTKlMGoHa7zMqoFDG9HQZt9vErrtHsj67Jx4NVo7YyRIKF6SOPCnfVz0NKXJaim6T0Eyfaa4xJGkDOSHM5QGvFw65n0yrek=
// term_id:88888888
// trans_stat:SUCCESS
// transaction_id:2022070122001476811446553236
// handerFuiouTradeQuery
func (res *CommonResponse) handerFuiouTradeQuery(content mxj.Map) mxj.Map {
	// 查询 交易结果标志：0：支付中请稍后查询，1：支付成功，2：支付失败，3：已撤销，4：撤销中请稍后查询，5：已全额退款，6：已部分退款，7：退款中请稍后查询
	data := mxj.New()
	data["status"] = WAITING // 状态
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		data["return_code"] = SUCCESS
		switch content["trans_stat"] {
		case "SUCCESS":
			data["status"] = SUCCESS
		case "REFUND":
			data["status"] = SUCCESS
		case "USERPAYING":
			data["status"] = USERPAYING
		case "PAYERROR":
			data["status"] = CLOSED
		case "CLOSED":
			data["status"] = CLOSED
		case "REVOKED":
			data["status"] = CLOSED
		case "NOTPAY":
			data["status"] = USERPAYING
		}
		total_amt, _ := strconv.ParseInt(content["order_amt"].(string), 10, 64)
		data["total_fee"] = total_amt
		if v, ok := content["reserved_fy_order_no"]; ok {
			data["bank_trade_no"] = v // 富有订单
		}
		data["trade_no"] = content["transaction_id"]
		if v, ok := content["mchnt_order_no"]; ok {
			data["out_trade_no"] = v
		}
		if v, ok := content["reserved_txn_fin_ts"]; ok && v != "" {
			data["time_end"] = v
		}
		// channel
		if v, ok := content["order_type"]; ok {
			switch v {
			case "WECHAT":
				data["wechat_open_id"] = content["buyer_id"]
			case "ALIPAY":
				data["alipay_open_id"] = content["buyer_id"]
			}
		}
	} else {
		data["return_code"] = "FAIL"
		data["status"] = CLOSED
		if content["result_code"] == "030010" {
			data["status"] = USERPAYING
		}
		if content["result_code"] == "010002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "9999" {
			data["status"] = WAITING
		}
		if content["result_code"] == "010001" {
			data["status"] = WAITING
		}
		if content["result_code"] == "2002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030003" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030004" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1013" {
			data["status"] = WAITING
		}
	}
	return data
}

// ins_cd:08A9999999
// mchnt_cd:0002900F0370542
// mchnt_order_no:513457061273811818895
// order_type:ALIPAY
// random_str:WOBJB5Q16KL9N0NWQGPOISBUWBEJORON
// refund_id: refund_order_no:513457061273811818895-1
// reserved_fy_settle_dt:20220703
// reserved_fy_trace_no:070008398190
// reserved_refund_amt:1
// result_code:030001
// result_msg:卖家余额不足
// sign:DgJKNx6UmabLhnep17ktFaTgwMWhmN4OlxGfceRX2MRoLYUpVPhhhTQx+ki8xJ7mm0CXnju55MJ8qiK1a74PVfDgmXGp9lLczFk2ANcs6jLGuBonIgLu93UZqT0vmAMmj6mcQO2Yaek1aaHVcUAfquUt+MjhGP9nWLWnvbtqH3U=
// term_id:88888888
// transaction_id:2022070122001476811446553236

// channel:fuiou content:map[ins_cd:08A9999999
// mchnt_cd:0002900F0370542
// mchnt_order_no:513457061273811818896
// order_type:ALIPAY
// random_str:ZACSC2SBMFQW1FTEL74XI280VU4GBWDT
// refund_id: refund_order_no:513457061273811818896-1
// reserved_fy_settle_dt:20220703
// reserved_fy_trace_no:070008398196
// reserved_promotion_detail:[{"amount":"0.01","type":"ALIPAYACCOUNT"}]
// reserved_refund_amt:1
// result_code:000000
// result_msg:SUCCESS
// sign:XCIszDvFDgx8fIOzZQ1u0c4MifaqAWpbTnrb6wpM26EPHuDqhq5rg83QsDBIK3riPuE/N69vRwTidZUmEYsGJS
// ZmWCPdJ10T/Qz9NG6mLnqnbYH81Pp/wv5AadKxFnjkcxYKgMoUWnxv8ibOM4ylfdK29RS6r3usCZtgHiIK22o=
// term_id:88888888
// transaction_id:2022070322001476811448079879
// ]
// out_refund_no:513457061273811818896-1
// out_trade_no:513457061273811818896
// refund_fee:1
// return_code:SUCCESS
// return_msg:SUCCESS
// status:SUCCESS
// trade_no:
// handerFuiouTradeRefund
func (res *CommonResponse) handerFuiouTradeRefund(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		data["return_code"] = SUCCESS
		data["status"] = SUCCESS
		switch content["result_msg"] {
		case "USERPAYING":
			data["status"] = USERPAYING
		case "CLOSED":
			data["status"] = CLOSED
		case "REVOKED":
			data["status"] = CLOSED
		case "NOTPAY":
			data["status"] = USERPAYING
		}
		if v, ok := content["reserved_refund_amt"]; ok {
			data["refund_fee"] = v
		}
		data["out_trade_no"] = content["mchnt_order_no"]
		data["out_refund_no"] = content["refund_order_no"]
		if v, ok := content["reserved_fy_settle_dt"]; ok && v != "" {
			data["time_end"] = v
		}
	} else {
		data["return_code"] = "FAIL"
		data["status"] = CLOSED
		if content["result_code"] == "030010" {
			data["status"] = USERPAYING
		}
		if content["result_code"] == "010002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "9999" {
			data["status"] = WAITING
		}
		if content["result_code"] == "010001" {
			data["status"] = WAITING
		}
		if content["result_code"] == "2002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030003" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030004" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1013" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1011" {
			data["status"] = USERPAYING
		}
	}
	return data
}

// ins_cd:08A9999999 mchnt_cd:0002900F0370542
// mchnt_order_no:513457061273811818896
// order_type:ALIPAY random_str:YB9KTTF2SBFNXPFK4SQE0ZPBTVSW95VP
// refund_id: refund_order_no:513457061273811818896-1
// reserved_fy_settle_dt:20220703
// reserved_fy_trace_no:070008398196
// reserved_promotion_detail:[{"amount":"0.01","type":"ALIPAYACCOUNT"}]
// reserved_refund_amt:1 result_code:000000
// result_msg:SUCCESS sign:XQ3LV5ZhF1ZHZTKqsvtNOUSulM/s7UGnICS4KxFCG/qGqvUqv8D8OTpn5m
// EHyvwFkpHmj9nJpa9ilqMM/clOQbOid40y0rFTuzUfL5Ei+Wo8OOg0b2IItGcS+p19oRlPhKHBMIwsxH1sHar
// 0g/1F4frI2z3/KutsgDVxtLoPhbI=
// term_id:88888888
// trans_stat:SUCCESS
// transaction_id:2022070322001476811448079879
// handerFuiouTradeRefundQuery
func (res *CommonResponse) handerFuiouTradeRefundQuery(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["status"] = WAITING // 状态
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		data["return_code"] = SUCCESS
		switch content["trans_stat"] {
		case "SUCCESS":
			data["status"] = SUCCESS
		case "USERPAYING":
			data["status"] = USERPAYING
		case "PAYERROR":
			data["status"] = CLOSED
		case "CLOSED":
			data["status"] = CLOSED
		case "REVOKED":
			data["status"] = CLOSED
		case "NOTPAY":
			data["status"] = USERPAYING
		}
		if v, ok := content["reserved_refund_amt"]; ok {
			data["refund_fee"] = v
		}
		data["out_trade_no"] = content["mchnt_order_no"]
		data["out_refund_no"] = content["refund_order_no"]
		if v, ok := content["reserved_fy_settle_dt"]; ok && v != "" {
			data["time_end"] = v
		}
	} else {
		data["return_code"] = "FAIL"
		data["status"] = CLOSED
		if content["result_code"] == "030010" {
			data["status"] = USERPAYING
		}
		if content["result_code"] == "010002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "9999" {
			data["status"] = WAITING
		}
		if content["result_code"] == "010001" {
			data["status"] = WAITING
		}
		if content["result_code"] == "2002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030003" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030004" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1013" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1011" {
			data["status"] = USERPAYING
		}
	}
	return data
}

func (res *CommonResponse) handerFuiouTradeQrcode(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		data["return_code"] = SUCCESS
		data["qr_code"] = content["qr_code"]
		data["status"] = USERPAYING
	} else {
		data["return_code"] = "FAIL"
		data["status"] = CLOSED
		if content["result_code"] == "030010" {
			data["status"] = USERPAYING
		}
		if content["result_code"] == "010002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "9999" {
			data["status"] = WAITING
		}
		if content["result_code"] == "010001" {
			data["status"] = WAITING
		}
		if content["result_code"] == "2002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030003" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030004" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1013" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1011" {
			data["status"] = USERPAYING
		}
	}
	return data
}

// handerFuiouTradeOpenid
func (res *CommonResponse) handerFuiouTradeOpenid(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		data["return_code"] = SUCCESS
		data["wechat_open_id"] = content["openid"]
		data["wechat_sub_open_id"] = content["sub_openid"]
		if v, ok := content["userId"]; ok && v != "" {
			data["alipay_user_id"] = v
		}
	} else {
		data["return_code"] = "FAIL"
	}
	return data
}

// ParseNotifyResult 解析异步通知
func (res *CommonResponse) InterfaceToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float32:
		return strconv.FormatFloat(v.(float64), 'E', -1, 32)
	case float64:
		return strconv.FormatFloat(v.(float64), 'E', -1, 64)
	}
	return ""
}

// HanderFuiouNotify
func (res *CommonResponse) HanderFuiouNotify(content mxj.Map) mxj.Map {
	// 查询 交易结果标志：0：支付中请稍后查询，1：支付成功，2：支付失败，3：已撤销，4：撤销中请稍后查询，5：已全额退款，6：已部分退款，7：退款中请稍后查询
	data := mxj.New()
	data["channel"] = "fuiou" //渠道
	data["content"] = content
	data["status"] = "" // 状态
	data["return_msg"] = content["return_msg"]
	if content["return_code"] == "0" {
		data["return_code"] = SUCCESS
		data["status"] = SUCCESS
		if v, ok := content["total_amt"]; ok {
			total_amt, _ := strconv.ParseInt(res.InterfaceToString(v), 10, 64)
			data["total_fee"] = total_amt
			if v, ok := content["payment_amt"]; ok { // 用户实际扣减金额
				payment_amt, _ := strconv.ParseInt(res.InterfaceToString(v), 10, 64)
				data["buyer_pay_fee"] = payment_amt
			} else {
				data["buyer_pay_fee"] = total_amt
			}
		}
		data["bank_trade_no"] = content["order_id"] // 银行订单
		data["trade_no"] = ""
		data["out_trade_no"] = content["out_trade_no"]
		data["time_end"] = content["pay_time"]
		if v, ok := content["channel"]; ok {
			switch v {
			case "91":
				data["method"] = "wechat"
				data["wechat_open_id"] = content["cust_id"]
			case "92":
				data["method"] = "alipay"
				data["alipay_logon_id"] = content["buyer_logon_id"]
				data["alipay_user_id"] = content["cust_id"]
			case "93":
				data["method"] = "unionpay"
			case "99":
				data["method"] = "unionpay"
			case "94":
				data["method"] = "digital"
			}
		}
	} else {
		data["return_code"] = "FAIL"
	}
	return data
}

// ins_cd:08A9999999
// mchnt_cd:0002900F0370542
// qr_code:
// random_str:BF6XB8YCHAOD4DY2HGVZ23YJ322JTVWU
// reserved_addn_inf:
// reserved_channel_order_id:
// reserved_fy_order_no: reserved_fy_settle_dt:
//  reserved_fy_trace_no:070008398303
//  reserved_pay_info:
//  reserved_transaction_id:
//  result_code:030001
//  result_msg:sub_mch_id与sub_appid不匹配
//  sdk_appid: sdk_noncestr: sdk_package:
//  sdk_partnerid: sdk_paysign:
//  sdk_signtype: sdk_timestamp:
//  session_id: sign:VkyE3g9m1R33x9dHr6bhanW07vghprI
//  vR5Zm6BdACnv3ebKzKapepZuRTR4DsaQGyO33Ep0ZuXZvIcX4Y4QOW2n/iFg1avjEkfNRU5jonWfs3o+I
//  eyLpG9HD8VbWXuujB0qPVOsxM5TVaVhRpYKQ8PaYmQO7XtXd6IVO7sQlvHU=
//  sub_appid:12211212 sub_mer_id:406875231 sub_openid:1212112
//  term_id:263575187
// handerFuiouTradeJsApi
func (res *CommonResponse) handerFuiouTradeJsApi(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["status"] = "" // 状态
	data["return_msg"] = content["result_msg"]
	if content["result_code"] == "000000" {
		// 交易结果标志,-1:下单失败，0：支付中，1：支付成功，2：支付失败
		data["return_code"] = SUCCESS
		data["status"] = USERPAYING
		data["time_end"] = content["pay_time"]
		if _, ok := content["sdk_appid"]; ok {
			mapWx := map[string]interface{}{
				"appId":     content["sdk_appid"],
				"nonceStr":  content["sdk_noncestr"],
				"package":   content["sdk_package"],
				"paySign":   content["sdk_paysign"],
				"signType":  content["sdk_signtype"],
				"timeStamp": content["sdk_timestamp"],
			}
			strWx, _ := json.Marshal(mapWx)
			data["wechat_package"] = string(strWx)
		}
		if v, ok := content["session_id"]; ok {
			data["prepay_id"] = v
		}
	} else {
		data["return_code"] = "FAIL"
		data["status"] = CLOSED
		if content["result_code"] == "030010" {
			data["status"] = USERPAYING
		}
		if content["result_code"] == "010002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "9999" {
			data["status"] = WAITING
		}
		if content["result_code"] == "010001" {
			data["status"] = WAITING
		}
		if content["result_code"] == "2002" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030003" {
			data["status"] = WAITING
		}
		if content["result_code"] == "030004" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1013" {
			data["status"] = WAITING
		}
		if content["result_code"] == "1011" {
			data["status"] = USERPAYING
		}
	}
	return data
}

// handerFuiouTradeB2cQuery
func (res *CommonResponse) handerFuiouTradeB2cQuery(content mxj.Map) mxj.Map {
	// 查询 交易结果标志：0：支付中请稍后查询，1：支付成功，2：支付失败，3：已撤销，4：撤销中请稍后查询，5：已全额退款，6：已部分退款，7：退款中请稍后查询
	data := mxj.New()
	data["status"] = WAITING // 状态
	data["return_msg"] = content["return_msg"]
	if content["return_code"] == "0" {
		data["return_code"] = SUCCESS
		switch content["pay_status"] {
		case "0":
			data["status"] = SUCCESS
		case "1":
			data["status"] = CLOSED
		case "2":
			data["status"] = WAITING
		}
		if v, ok := content["total_amt"]; ok {
			total_amt, _ := strconv.ParseInt(v.(string), 10, 64)
			data["total_fee"] = total_amt
			if v, ok := content["payment_amt"]; ok { // 用户实际扣减金额
				i, _ := strconv.ParseInt(v.(string), 10, 64)
				data["buyer_pay_fee"] = i
			} else {
				data["buyer_pay_fee"] = total_amt
			}
		}
		data["bank_trade_no"] = content["order_id"] // 银行订单
		data["trade_no"] = content["third_trade_no"]
		data["out_trade_no"] = content["out_trade_no"]
		data["time_end"] = content["pay_time"]
		if v, ok := content["pay_type"]; ok {
			switch v {
			case "9":
				data["method"] = "wechat"
				data["wechat_open_id"] = content["open_id"]
			case "10":
				data["method"] = "alipay"
				data["alipay_user_id"] = content["open_id"]
			default:
				data["method"] = "unionpay"
			}
		}
	} else {
		data["return_code"] = "FAIL"
		if content["return_code"] == "20000002" {
			data["return_code"] = SUCCESS
			switch content["pay_status"] {
			case "0":
				data["status"] = SUCCESS
			case "1":
				data["status"] = CLOSED
			case "2":
				data["status"] = WAITING
			}
		}
		// "return_code": "12001081",
		// "return_msg": "查询订单对照表获取工行订单号失败",
		if content["return_code"] == "12001081" {
			data["return_code"] = SUCCESS
			data["status"] = CLOSED
		}
		// "return_code": "12001081",
		// "return_msg": "查询订单对照表获取工行订单号失败",
		if content["return_code"] == "12001081" {
			data["return_code"] = SUCCESS
			data["status"] = CLOSED
		}
		if res.InterfaceToString(content["return_code"]) == "400019" {
			data["status"] = CLOSED
		}
	}
	return data
}

// {
// 	"response_biz_content": {
// 		"ecoupon_amt": "0",
// 		"access_type": "8",
// 		"card_kind": "",
// 		"open_id": "2088002104076813",
// 		"return_msg": "SUCCESS",
// 		"third_party_coupon_amt": "0",
// 		"total_amt": "1",
// 	"pay_status": "0",
// 	"card_no": "",
// 	"bank_disc_amt": "0",
// 	"third_party_discount_amt": "0",
// 	"decr_flag": "",
// 	"pay_type": "10",
// 	"mer_disc_amt": "0",
// 	"attach": "0de12a77-3f0a-44fc-b175-64eaee244bf3",
// 	"msg_id": "160304042047000532111040002054",
// 	"point_amt": "0",
// 	"third_trade_no": "802021110422001476815755339317",
// 	"mer_id": "160304042047",
// 	"card_flag": "",
// 	"total_disc_amt": "0",
// 	"pay_time": "20211104164117",
// 	"out_trade_no": "20211104164111482994",
// 	"coupon_amt": "0",
// 	"payment_amt": "1",
// 	"return_code": "0",
// 	"order_id": "160304042047000532111040002054"
// }

// {
// 	"response_biz_content": {
// 		"access_type": "",
// 		"pay_status": "2",
// 		"card_kind": "",
// 		"third_party_discount_amt": "0",
// 		"card_flag": "",
// 		"open_id": "",
// 		"decr_flag": "",
// 		"return_msg": "调SAES_AC通知查询关单ACS服务未知失败通知查询关单接口处理失败:第三方交易状态异常：0，第三方订单待支付",
// 		"third_party_coupon_amt": "0",
// 		"pay_type": "0",
// 		"return_code": "20000002",
// 		"order_id": "160304042047000542111040002057"
// 	}

// handerFuiouTradeB2cRefund
func (res *CommonResponse) handerFuiouTradeB2cRefund(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["return_msg"] = content["return_msg"]
	if content["return_code"] == "0" {
		data["return_code"] = SUCCESS
		data["status"] = SUCCESS
		// data["total_fee"] = content["total_fee"]
		reject_amt, _ := strconv.ParseInt(content["reject_amt"].(string), 10, 64)
		data["refund_fee"] = reject_amt
		data["bank_trade_no"] = content["order_id"]
		data["out_trade_no"] = content["out_trade_no"]
		data["out_refund_no"] = content["outtrx_serial_no"]
	} else {
		data["return_code"] = "FAIL"
		if content["return_code"] == "00070038" {
			data["status"] = WAITING
		}
		if res.InterfaceToString(content["return_code"]) == "400019" {
			data["status"] = CLOSED
		}
	}
	return data
}

// {"response_biz_content":{"reject_mer_disc_amt":"0",
// "reject_ecoupon":"0","settlement_refund_amt":"1",
// "intrx_serial_no":"160304042047000532111040003110",
// "return_msg":"成功","third_party_return_msg":"",
// "reject_amt":"1","reject_bank_disc_amt":"0",
// "reject_point":"0","out_trade_no":"20211104164111482994",
// "card_no":"","third_party_discount_refund_amt":"0",
// "third_party_coupon_refund_amt":"0",
// "third_party_return_code":"","real_reject_amt":"1",
// "pay_type":"10","return_code":"0",
// "order_id":"160304042047000532111040002054",
// "outtrx_serial_no":"T20211104164111482994"}

// handerFuiouTradeB2cRefundQuery
func (res *CommonResponse) handerFuiouTradeB2cRefundQuery(content mxj.Map) mxj.Map {
	data := mxj.New()
	data["return_msg"] = content["return_msg"]
	if content["return_code"] == "0" {
		data["return_code"] = SUCCESS
		switch content["pay_status"] {
		case "0":
			data["status"] = SUCCESS
		case "1":
			data["status"] = CLOSED
		case "2":
			data["status"] = WAITING
		}
		// data["total_fee"] = content["total_fee"]
		reject_amt, _ := strconv.ParseInt(content["reject_amt"].(string), 10, 64)
		data["refund_fee"] = reject_amt
		data["bank_trade_no"] = content["order_id"]
		data["out_trade_no"] = content["out_trade_no"]
		data["out_refund_no"] = content["outtrx_serial_no"]
	} else {
		data["return_code"] = "FAIL"
		switch content["pay_status"] {
		case "0":
			data["status"] = SUCCESS
		case "1":
			data["status"] = CLOSED
		case "2":
			data["status"] = WAITING
		}
		if res.InterfaceToString(content["return_code"]) == "400019" {
			data["status"] = CLOSED
		}
	}
	return data
}
