/**
 * @Author: YMBoom
 * @Description:
 * @File:  util
 * @Version: 1.0.0
 * @Date: 2022/12/05 19:21
 */
package utilx

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func GetSign(timestamp, nonce, key, secret string) string {
	str := fmt.Sprintf("clientKey=%s&clientSecret=%s&nonce=%s&timestamp=%s", key, secret, nonce, timestamp)
	return strings.ToUpper(MD5(str))
}
func GetNonce() string {
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

func BaseHeader(token string, refererNum int) map[string]string {
	return map[string]string{
		"Authorization": token,
		"content-type":  "application/json",
		"User-Agent":    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.30(0x18001e31) NetType/WIFI Language/zh_CN",
		"Referer":       fmt.Sprintf("https://servicewechat.com/wx8e45b2134cbeddff/%d/page-frame.html", refererNum),
	}
}

func HttpGet(url string, header map[string]string) (res []byte, err error) {
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

// WriteLog 写入日志
func WriteLog(path, data string) {
	data = fmt.Sprintf("时间：%s ==== 内容：%s\r", time.Now().Format("2006-01-02 15:04:05.000000000"), data)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	defer file.Close()

	file.WriteString(data)
}
