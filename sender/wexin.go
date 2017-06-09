package sender

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/patrickmn/go-cache"
	"github.com/sdvdxl/falcon-message/util"
	"github.com/sdvdxl/go-tools/encrypt"
)

type Weixin struct {
	tokenCache     *cache.Cache
	CorpID         string
	Secret         string
	AgentID        string
	EncodingAESKey string
}

func (wx Weixin) Auth(echostr string) ([]byte, error) {
	wByte, err := base64.StdEncoding.DecodeString(echostr)
	if err != nil {
		return nil, errors.New("接受微信请求参数 echostr base64解码失败(" + err.Error() + ")")
	}
	key, err := base64.StdEncoding.DecodeString(wx.EncodingAESKey + "=")
	if err != nil {
		return nil, errors.New("配置 EncodingAESKey base64解码失败(" + err.Error() + "), 请检查配置文件内 EncodingAESKey 是否和微信后台提供一致")
	}

	keyByte := []byte(key)
	x := encrypt.AesDecrypt(wByte, keyByte)

	buf := bytes.NewBuffer(x[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)

	//验证返回数据ID是否正确
	appIDstart := 20 + length
	if len(x) < int(appIDstart) {
		return nil, errors.New("获取数据错误, 请检查 EncodingAESKey 配置")
	}
	id := x[appIDstart : int(appIDstart)+len(wx.CorpID)]
	if string(id) == wx.CorpID {
		return x[20 : 20+length], nil
	}
	return nil, errors.New("微信验证appID错误, 微信请求值: " + string(id) + ", 配置文件内配置为: " + wx.CorpID)
}

func NewWeixin(corpId, secret string) *Weixin {
	if corpId == "" || secret == "" {
		log.Fatal("corpId或者secret 获取失败, 请检查配置文件")
	}
	return &Weixin{tokenCache: cache.New(6000*time.Second, 5*time.Second)}
}

//发送信息
type content struct {
	Content string `json:"content"`
}

type msgPost struct {
	ToUser  string  `json:"touser"`
	MsgType string  `json:"msgtype"`
	AgentID int     `json:"agentid"`
	Text    content `json:"text"`
}

type accessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

const (
	weixinURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid="
)

// GetAccessToken 从微信获取 AccessToken
func (wx Weixin) GetAccessToken() {
	for {

		wxAccessTokenRUL := weixinURL + wx.CorpID + "&corpsecret=" + wx.Secret

		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}

		client := &http.Client{Transport: tr}
		result, err := client.Get(wxAccessTokenRUL)
		if err != nil {
			log.Printf("获取微信 Token 返回数据错误: %v", err)
			return
		}

		res, err := ioutil.ReadAll(result.Body)

		if err != nil {
			log.Printf("获取微信 Token 返回数据错误: %v", err)
			return
		}
		newAccess := accessToken{}
		err = json.Unmarshal(res, &newAccess)
		if err != nil {
			log.Printf("获取微信 Token 返回数据解析 Json 错误: %v", err)
			return
		}

		if newAccess.ExpiresIn == 0 || newAccess.AccessToken == "" {
			log.Printf("获取微信错误代码: %v, 错误信息: %v", newAccess.ErrCode, newAccess.ErrMsg)
			time.Sleep(5 * time.Minute)
		}

		//延迟时间
		wx.tokenCache.Set("token", newAccess, time.Duration(newAccess.ExpiresIn)*time.Second)
		log.Printf("微信 Token 更新成功: %s,有效时间: %v", newAccess.AccessToken, newAccess.ExpiresIn)
		time.Sleep(time.Duration(newAccess.ExpiresIn-100) * time.Second)
	}

}

// WxPost 微信请求数据
func wxPost(url string, data msgPost) (string, error) {
	jsonBody, err := util.EncodeJSON(data)
	if err != nil {
		return "", err
	}

	r, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func (wx Weixin) Send(tos string, message string) error {
	var toUser string
	if userList := strings.Split(tos, ","); len(userList) > 1 {
		toUser = strings.Join(userList, "|")
	}

	text := content{}
	text.Content = message

	msg := msgPost{
		ToUser:  toUser,
		MsgType: "text",
		AgentID: util.StringToInt(wx.AgentID),
		Text:    text,
	}

	token, found := wx.tokenCache.Get("token")
	if !found {
		return errors.New("token获取失败")
	}
	accessToken, ok := token.(accessToken)
	if !ok {
		return errors.New("token解析失败")
	}

	url := weixinURL + accessToken.AccessToken

	result, err := wxPost(url, msg)
	if err != nil {
		return err
	}
	log.Printf("发送信息给%s, 信息内容: %s, 微信返回结果: %v", toUser, message, result)
	return nil
}
