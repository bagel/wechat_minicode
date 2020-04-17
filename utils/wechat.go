package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// WechatCode struct
type WechatCode struct {
	AppID       string
	AppSecret   string
	Scene       string
	Width       int
	AccessToken string
}

// WechatAccessTokenResponse access token response
type WechatAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// WechatCodeRequest request wx code
type WechatCodeRequest struct {
	Scene     string `json:"scene"`
	Page      string `json:"page"`
	Width     int    `json:"width"`
	AutoColor bool   `json:"auto_color"`
}

// NewWechatCode init
func NewWechatCode(appID string, appSecret string, scene string, width int) *WechatCode {
	return &WechatCode{
		AppID:     appID,
		AppSecret: appSecret,
		Scene:     scene,
		Width:     width,
	}
}

// GetWechatAccessToken get access token
// Wechat doc https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
func (w *WechatCode) GetWechatAccessToken() (err error) {
	v := url.Values{}
	v.Set("grant_type", "client_credential")
	v.Set("appid", w.AppID)
	v.Set("secret", w.AppSecret)

	wxURL := "https://api.weixin.qq.com/cgi-bin/token?" + v.Encode()
	log.Println(wxURL)

	resp, err := http.Get(wxURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var responseBody WechatAccessTokenResponse
	if err = json.Unmarshal([]byte(body), &responseBody); err != nil {
		return
	}

	//log.Println(responseBody.AccessToken)

	w.AccessToken = responseBody.AccessToken

	return
}

// GetWechatCode wechat code
// Wechat doc https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func (w *WechatCode) GetWechatCode() (body []byte, err error) {
	if err = w.GetWechatAccessToken(); err != nil {
		return
	}

	wxURL := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + w.AccessToken

	request := WechatCodeRequest{
		Scene:     w.Scene,
		Page:      "pages/index/index",
		Width:     w.Width,
		AutoColor: true,
	}

	postBody, err := json.Marshal(request)
	if err != nil {
		return
	}

	resp, err := http.Post(wxURL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if _, ok := resp.Header["Content-Type"]; !ok {
		err = errors.New("Wechat reponse header Content-Type not found")
		return
	}

	contentType := resp.Header["Content-Type"][0]
	if !strings.Contains(contentType, "image") {
		err = errors.New("Wechat response error")
		return
	}

	return
}
