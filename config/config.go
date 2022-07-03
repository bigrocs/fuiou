package config

type Config struct {
	InsCd          string `json:"ins_cd"`           //工商银行分配给开发者的应用ID
	PrivateKey     string `json:"private_key"`      //私钥
	FuiouPublicKey string `json:"fuiou_public_key"` //工商银行公钥
	MchntCd        string `json:"mchnt_cd"`         //消息通讯唯一编号，每次调用独立生成，APP级唯一
	Version        string `json:"version"`          //仅支持 JSON
	TermId         string `json:"term_id"`          //请求使用的编码格式，如utf-8,gbk,gb2312等，推荐使用 utf-8
	Sign           string `json:"sign"`             //商户请求参数的签名串
	NotifyUrl      string `json:"notify_url"`       //工商银行服务器主动通知商户服务器里指定的页面http/https路径。
	BizContent     string `json:"biz_content"`      //业务请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
	WechatAppId    string `json:"wechat_app_id"`    // 微信支付时绑定的微信公众号
	Sandbox        bool   `json:"sandbox"`          // 沙盒
}
