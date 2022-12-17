/**
 * @Author: YMBoom
 * @Description:
 * @File:  resp
 * @Version: 1.0.0
 * @Date: 2022/12/17 16:34
 */
package seckill

type secResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data"`
}

type getJsonIdResp struct {
	ActivityData []struct {
		ResultId  string `json:"result_id"`
		JsonId    string `json:"json_id"`
		RewardNum int64  `json:"reward_num"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
	} `json:"activity_data"`
}

type getRkResp struct {
	Code int `json:"code"`
	Data struct {
		Rk string `json:"rk"`
	} `json:"data"`
}
