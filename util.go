/**
 * @Author: YMBoom
 * @Description:
 * @File:  util
 * @Version: 1.0.0
 * @Date: 2022/12/05 19:21
 */
package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func getTimestamp() string {
	//fmt.Println(time.Now().UnixNano())
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func getSign(timestamp, nonce string) string {
	clientKey := "l5lEyVkz3rXfDtCH41FAQKVueoX9HM4ZXmtNZn8ptajm88KWXMViVvcBrZLva9KO"
	clientSecret := "g3nbZazVhVOYaoLGD5mJAcKOLMQ28v7kh1KBHdekNkA9a43txBC2DBAPQtf0JJYm6FjdJD9AnwJELhRQS8F7Aia0yzAbpd9SqohiZBtLbbeYed1ada83sgJzZad49uC2"

	str := fmt.Sprintf("clientKey=%s&clientSecret=%s&nonce=%s&timestamp=%s", clientKey, clientSecret, nonce, timestamp)
	return strings.ToUpper(MD5(str))
}
func getNonce() string {
	return RandStr(16)
}

// RandStr 生成字母随机字符串
func RandStr(l int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"
	b := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}

func MD5(v string) string {
	m := md5.New()
	m.Write([]byte(v))

	x := m.Sum(nil)
	return fmt.Sprintf("%x", x)
}

func baseHeader() map[string]string {
	return map[string]string{
		"Authorization": "HT8E9Q8yr19iAq5/sgktychoHXxP3Yfzhlgep39iijyb+e+oILoMqctiHXUyQEya",
		"content-type":  "application/json",
		"User-Agent":    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.30(0x18001e31) NetType/WIFI Language/zh_CN",
		"Referer":       "https://servicewechat.com/wx8e45b2134cbeddff/56/page-frame.html",
	}
}

func httpGet(url string, header map[string]string) (res []byte, err error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	if header != nil {
		for i, m := range header {
			req.Header.Set(i, m)
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("do err", err)
		return nil, err
	}

	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return all, nil
}
