package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/djimenez/iconv-go"
	uuid "github.com/satori/go.uuid"

	"github.com/bigrocs/fuiou/config"
	"github.com/bigrocs/fuiou/requests"
	"github.com/bigrocs/fuiou/responses"
	"github.com/bigrocs/fuiou/util"
	"github.com/micro/go-micro/v2/util/log"
)

var apiUrlsMch = map[string]string{
	"pay.pay":           "/micropay",      //付款码支付
	"pay.orderquery":    "/commonQuery",   //付款码支付查询
	"pay.hisorderquery": "/hisTradeQuery", //付款码支付查询
	"pay.refund":        "/commonRefund",  //退款
	"pay.refundquery":   "/refundQuery",   //退款查询

	"pay.qrcode": "/preCreate",   //生成二维码
	"pay.jsapi":  "/wxPreCreate", //聚合支付第三方用户标识查询
	"pay.openid": "/auth2Openid", //聚合支付第三方用户标识查询

}

// Common 公共封装
type Common struct {
	Config   *config.Config
	Requests *requests.CommonRequest
}

// Action 创建新的公共连接
func (c *Common) Action(response *responses.CommonResponse) (err error) {
	return c.Request(response)
}

// APIBaseURL 默认 API 网关
func (c *Common) APIBaseURL() string { // TODO(): 后期做容灾功能
	con := c.Config
	if con.Sandbox { // 沙盒模式
		return "https://fundwx.fuiou.com"
	}
	return "https://fundwx.fuiou.com"
}

// ApiUrl 创建 ApiUrl
func (c *Common) ApiUrl() (apiUrl string, err error) {
	req := c.Requests
	if u, ok := apiUrlsMch[req.ApiName]; ok {
		apiUrl = c.APIBaseURL() + u
	} else {
		err = fmt.Errorf("ApiName 不存在请检查。")
	}
	return
}

// Request 执行请求
// AppId        string `json:"app_id"`         //工行开发平台分配给开发者的应用ID
// Method       string `json:"method"`         //接口名称
// Format       string `json:"format"`         //仅支持 JSON
// Charset      string `json:"charset"`        //请求使用的编码格式，如utf-8,gbk,gb2312等，推荐使用 utf-8
// SignType     string `json:"sign_type"`      //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用 RSA2
// Sign         string `json:"sign"`           //商户请求参数的签名串
// Timestamp    string `json:"timestamp"`      //发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
// Version      string `json:"version"`        //调用的接口版本，固定为：1.0
// NotifyUrl    string `json:"notify_url"`     //工行开发平台服务器主动通知商户服务器里指定的页面http/https路径。
// BizContent   string `json:"biz_content"`    //业务请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
// ReturnUrl    string `json:"return_url"`     //HTTP/HTTPS开头字符串
func (c *Common) Request(response *responses.CommonResponse) (err error) {
	con := c.Config
	req := c.Requests
	apiUrl, err := c.ApiUrl()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	// 构建配置参数
	params := map[string]interface{}{
		"version":    "1.0",
		"ins_cd":     con.InsCd,
		"mchnt_cd":   con.MchntCd,
		"term_id":    "88888888",
		"random_str": strings.Replace(uuid.NewV4().String(), "-", "", -1), // 32位
	}
	if con.Version != "" {
		params["version"] = con.Version
	} else {
		params["version"] = "1.0"
	}
	if con.TermId != "" {
		params["term_id"] = con.TermId
	} else {
		params["term_id"] = "88888888"
	}
	if req.ApiName == "pay.openid" {
		// 删除params中的term_id字段
		delete(params, "term_id")
	}
	// 整合 params
	for k, v := range req.BizContent {
		// 包含中文的 goods_des 字段转gbk然后encode
		if k == "goods_des" {
			output, _ := iconv.ConvertString(v.(string), "utf-8", "gbk")
			v = output
			req.BizContent[k] = v
		}
		if k == "goods_detail" {
			output, _ := iconv.ConvertString(v.(string), "utf-8", "gbk")
			v = output
			req.BizContent[k] = v
		}
		if k == "goods_tag" {
			output, _ := iconv.ConvertString(v.(string), "utf-8", "gbk")
			v = output
			req.BizContent[k] = v
		}
		if k == "addn_inf" {
			output, _ := iconv.ConvertString(v.(string), "utf-8", "gbk")
			v = output
			req.BizContent[k] = v
		}
		// 检测k字符串是否存在reserved
		if !strings.HasPrefix(k, "reserved") {
			params[k] = v
		}
	}
	sign, err := util.Sign(params, con.PrivateKey, "MD5") // 开发签名
	if err != nil {
		return err
	}
	params["sign"] = sign
	// 整合 params 的reserved字段
	for k, v := range req.BizContent {
		// 检测k字符串是否存在reserved
		if strings.HasPrefix(k, "reserved") {
			params[k] = v
		}
	}
	obj, err := util.GbkXML(params)
	if err != nil {
		return err
	}
	log.Info("Fuiou[PostForm]", apiUrl, obj)
	res, err := util.PostForm(apiUrl, `req=`+url.QueryEscape(url.QueryEscape(obj)))
	if err != nil {
		log.Info("Fuiou[PostForm]err", err)
		return err
	}
	response.SetHttpContent(res, "xml")
	return
}
