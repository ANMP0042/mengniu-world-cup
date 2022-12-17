/**
 * @Author: YMBoom
 * @Description:
 * @File:  base_info
 * @Version: 1.0.0
 * @Date: 2022/12/07 9:48
 */
package seckill

import (
	"encoding/json"
	"fmt"
	"mengniu/utilx"
)

type RkResp struct {
	Code int `json:"code"`
	Data struct {
		Rk string `json:"rk"`
	} `json:"data"`
}

func getRk(token string) string {
	timestamp := utilx.GetTimestamp()
	nonce := utilx.GetNonce()
	sign := utilx.GetSign(timestamp, nonce)
	url := fmt.Sprintf("%s/%s?timestamp=%s&nonce=%s&signature=%s", utilx.Domain, "user/baseInfo", timestamp, nonce, sign)

	res, err := utilx.HttpGet(url, utilx.BaseHeader(token))
	if err != nil {
		fmt.Println("获取rk失败，err= ", err)
		return ""
	}

	fmt.Println("rk resp ", string(res))
	rk := new(RkResp)
	if err = json.Unmarshal(res, rk); err != nil {
		fmt.Println("rkResp 解析失败，err= ", err)
		return ""
	}

	return rk.Data.Rk
}
