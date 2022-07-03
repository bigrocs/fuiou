package util

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/djimenez/iconv-go"
)

const (
	SignType_MD5    = "MD5"
	SignType_SHA1   = "SHA1"
	SignType_SHA256 = "SHA256"
)

// VerifySign 验证支付
func VerifySign(params string, sign string, fuiouPublicKey string, signType string) (ok bool, err error) {
	p, err := mxj.NewMapJson([]byte(params))
	if err != nil {
		return false, err
	}
	signData := EncodeSignParams(p)
	signData, err = iconv.ConvertString(signData, "utf-8", "gbk")
	if err != nil {
		return false, err
	}
	var (
		h         hash.Hash
		hashs     crypto.Hash
		block     *pem.Block
		pubKey    interface{}
		publicKey *rsa.PublicKey
		keyOk     bool
	)
	signBytes, _ := base64.StdEncoding.DecodeString(sign)
	if block, _ = pem.Decode([]byte(FormatPrivateKey(fuiouPublicKey))); block == nil {
		return ok, errors.New("富友开发平台公钥Decode错误")
	}
	if pubKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return ok, fmt.Errorf("x509.ParsePKIXPublicKey：%w", err)
	}
	if publicKey, keyOk = pubKey.(*rsa.PublicKey); !keyOk {
		return ok, errors.New("富友开发平台公钥转换错误")
	}
	switch signType {
	case "MD5":
		hashs = crypto.MD5
	case "RSA":
		hashs = crypto.SHA1
	case "RSA2":
		hashs = crypto.SHA256
	default:
		hashs = crypto.SHA256
	}
	h = hashs.New()
	h.Write([]byte(signData))
	err = rsa.VerifyPKCS1v15(publicKey, hashs, h.Sum(nil), signBytes)
	if err != nil {
		return ok, err
	}
	return true, err
}

// EncodeSignParams 编码符号参数
func EncodeSignParams(params map[string]interface{}) string {
	var buf strings.Builder
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		if strings.HasPrefix(k, "reserved") {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := params[k]
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(InterfaceToString(v))
		buf.WriteByte('&')
	}
	return buf.String()[:buf.Len()-1]
}

// Sign 工行开发平台签名支付签名.
//  params: 待签名的参数集合
//  privateKey: 密钥
func Sign(params map[string]interface{}, privateKey string, signType string) (sign string, err error) {
	encodeSignParams := EncodeSignParams(params)
	var (
		block          *pem.Block
		h              hash.Hash
		key            *rsa.PrivateKey
		hashs          crypto.Hash
		encryptedBytes []byte
	)
	if block, _ = pem.Decode([]byte(FormatPrivateKey(privateKey))); block == nil {
		return "", errors.New("pem.Decode：privateKey decode error")
	}
	itf, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	key, ok := itf.(*rsa.PrivateKey)
	if !ok {
		err = errors.New("fuiou: unknown type of public key")
	}
	switch signType {
	case "MD5":
		h = md5.New()
		hashs = crypto.MD5
	case "RSA":
		h = sha1.New()
		hashs = crypto.SHA1
	case "RSA2":
		h = sha256.New()
		hashs = crypto.SHA256
	default:
		h = sha256.New()
		hashs = crypto.SHA256
	}
	if _, err = h.Write([]byte(encodeSignParams)); err != nil {
		return
	}
	if encryptedBytes, err = rsa.SignPKCS1v15(rand.Reader, key, hashs, h.Sum(nil)); err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(encryptedBytes)
	return
}

// ParseNotifyResult 解析异步通知
func InterfaceToString(v interface{}) string {
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

// FormatPrivateKey 格式化 普通应用秘钥
func FormatPrivateKey(privateKey string) (pKey string) {
	var buffer strings.Builder
	buffer.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	rawLen := 64
	keyLen := len(privateKey)
	raws := keyLen / rawLen
	temp := keyLen % rawLen
	if temp > 0 {
		raws++
	}
	start := 0
	end := start + rawLen
	for i := 0; i < raws; i++ {
		if i == raws-1 {
			buffer.WriteString(privateKey[start:])
		} else {
			buffer.WriteString(privateKey[start:end])
		}
		buffer.WriteByte('\n')
		start += rawLen
		end = start + rawLen
	}
	buffer.WriteString("-----END RSA PRIVATE KEY-----\n")
	pKey = buffer.String()
	return
}

// FormatURLParam 格式化请求URL参数
func FormatURLParam(params map[string]interface{}) (urlParam string) {
	v := url.Values{}
	for key, value := range params {
		v.Add(key, InterfaceToString(value))
	}
	return v.Encode()
}

// getSignData 获取数据字符串
func GetSignData(str string) (signData string) {
	indexStart := strings.Index(str, `":`)
	indexEnd := strings.Index(str, `,"sign"`)
	if indexEnd == -1 {
		indexEnd = strings.Index(str, `}}`) + 1
	}
	signData = str[indexStart+2 : indexEnd]
	return
}
