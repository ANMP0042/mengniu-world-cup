/**
 * @Author: YMBoom
 * @Description:
 * @File:  seckill
 * @Version: 1.0.0
 * @Date: 2022/12/07 9:49
 */
package seckill

import (
	"encoding/json"
	"errors"
	"fmt"
	"mengniu/config"
	"mengniu/utilx"
	"strings"
	"sync"
	"time"
)

type (
	Seckiller struct {
		wg        sync.WaitGroup
		tokens    []tokens
		cnf       *config.Config
		jsonId    string
		rewardNum int64
		writeLog  bool
		writePath string
		ts        chan tokens
	}

	tokens struct {
		Token string
		Rk    string
	}

	Option func(reg *Seckiller)
)

func NewSeckiller(opts ...Option) (*Seckiller, error) {
	// 加载配置
	cnf, err := config.Load()

	if err != nil {
		return nil, errors.New("配置加载失败：" + err.Error())
	}

	seckiller := &Seckiller{
		cnf: cnf,
		ts:  make(chan tokens),
	}

	for _, opt := range opts {
		opt(seckiller)
	}
	return seckiller, nil
}
func (s *Seckiller) readToken() {

	for {
		ts := <-s.ts
		fmt.Println("收到一个token需要抢购", ts)
	}

}
func (s *Seckiller) Seckill() {
	if err := s.verify(); err != nil {
		fmt.Println("抢购发生错误:" + err.Error())
		return
	}

	// 设置定时器
	t, err := s.duration()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	<-t

	// 多账号同时抢购
	s.wg.Add(len(s.tokens))
	for _, ts := range s.tokens {
		go s.seckill(&ts)
	}

	s.wg.Wait()
	s.log("================ 抢购结束 10分钟后查询结果 ================")
	// 定时10分钟后查询结果
	//go func() {
	//	resultT := time.After(10 * time.Minute)
	//
	//	<-resultT
	//	s.Result()
	//}()

	return
}

func (s *Seckiller) verify() error {
	if len(s.tokens) == 0 {
		return errors.New("所有token获取rk失败，请检查token或clientKey等是否有效")
	}

	if s.jsonId == "" {
		return errors.New("jsonId 获取失败")
	}

	return nil
}

// 定时器时间
func (s *Seckiller) duration() (<-chan time.Time, error) {
	seckillAt := fmt.Sprintf("%s %s:%s", time.Now().Format("2006-01-02"), s.cnf.Sec.At, "00")
	parse, _ := time.ParseInLocation("2006-01-02 15:04:05", seckillAt, time.Local)

	// 如果时间过了 就定时明天 更换jsonId
	if parse.Sub(time.Now()) < 0 {
		parse = parse.AddDate(0, 0, 1)
		go func() {
			s.jsonId, s.rewardNum = s.getJsonId(1)
		}()
	}

	diff := parse.Sub(time.Now())

	if diff <= 0 {
		return nil, errors.New("时间过了")
	}

	diff = diff - time.Duration(s.cnf.Custom.PreDuration)*time.Millisecond

	diff = 1 * time.Second
	go fmt.Println("距离开抢还有：", diff, " === 今日抢购数量：", s.rewardNum)
	go s.log(fmt.Sprintf("距离开抢还有：%d === 今日抢购数量：%d", diff, s.rewardNum))

	// 定时器
	return time.After(diff), nil
}

/**
抢购逻辑判断 需要判断返回状态
1   排队，大概率中奖，不需要在抢购不然会返回-1
-1  从排队队伍中清除
429 请求过快
500 发生错误

*/
func (s *Seckiller) seckill(ts *tokens) {
	defer s.wg.Done()
	for i := 1; i <= s.cnf.Custom.SecTime; i++ {
		switch s.do(ts) {
		case 1:
			s.log(fmt.Sprintf("排队成功 当前token %s 已停止抢购，抢购次数：%d", ts.Token, i))
			return
		case -1:
			s.log(fmt.Sprintf("当前token %s 排队取消，停止抢购，抢购次数：%d", ts.Token, i))
			return
		case 500:
			s.log(fmt.Sprintf("系统错误，当前token %s 停止抢购，抢购次数：%d", ts.Token, i))
			return
		default:
			go fmt.Println(fmt.Sprintf("抢购次数：%d，继续执行", i))
		}
	}
}

// 执行抢购
func (s *Seckiller) do(ts *tokens) int {

	timestamp := utilx.GetTimestamp()
	nonce := utilx.GetNonce()

	rkSign := utilx.GetSign(timestamp, nonce, s.cnf.Sec.ClientKey, s.cnf.Sec.ClientSecret)
	sign, reqId := seckillSign(ts.Rk, timestamp, s.cnf.Sec.DesKey)
	u := fmt.Sprintf("%s%s?timestamp=%s&nonce=%s&signature=%s&jsonId=%s", s.cnf.Sec.Domain1122, s.cnf.Sec.Path, timestamp, nonce, rkSign, s.jsonId)

	header := utilx.BaseHeader(ts.Token, s.cnf.Sec.RefererNum)
	header["sign"] = sign
	header["timestamp"] = timestamp
	header["requestId"] = reqId
	res, err := utilx.HttpGet(u, header)
	if err != nil {
		s.log(fmt.Sprintf("抢购错误:%s,res:%s", err.Error(), string(res)))
		return 500
	}

	s.log(fmt.Sprintf("token：%s,res：%s", ts.Token, string(res)))
	resp := new(secResp)
	if err = json.Unmarshal(res, resp); err != nil {
		return 500
	}

	return resp.Code
}

func (s *Seckiller) delToken(ts tokens) {
	for i := 0; i < len(s.tokens); i++ {
		if s.tokens[i].Token == ts.Token {
			s.tokens = append(s.tokens[:i], s.tokens[i+1:]...)
			i--
		}
	}
}
func (s *Seckiller) log(data string) {
	if s.writeLog {
		go utilx.WriteLog(s.writePath, data)
	}
}

func (s *Seckiller) getJsonId(addDay int) (jsonId string, rewardNum int64) {
	res, err := utilx.HttpGet(s.cnf.Sec.JsonIdUrl, utilx.BaseHeader("", s.cnf.Sec.RefererNum))
	if err != nil {
		fmt.Println("获取jsonId失败，err= ", err)
		return
	}

	resp := new(getJsonIdResp)
	if err = json.Unmarshal(res, resp); err != nil {
		fmt.Println("rkResp 解析失败，err= ", err)
		return
	}

	now := time.Now().AddDate(0, 0, addDay)
	for _, data := range resp.ActivityData {
		t := time.UnixMilli(data.StartTime)
		if t.Day() == now.Day() {
			return data.JsonId, data.RewardNum
		}
	}
	return
}

func (s *Seckiller) getRk(token string) string {
	timestamp := utilx.GetTimestamp()
	nonce := utilx.GetNonce()
	sign := utilx.GetSign(timestamp, nonce, s.cnf.Sec.ClientKey, s.cnf.Sec.ClientSecret)
	url := fmt.Sprintf("%s/%s?timestamp=%s&nonce=%s&signature=%s", s.cnf.Sec.Domain, "mp/api/user/baseInfo", timestamp, nonce, sign)

	res, err := utilx.HttpGet(url, utilx.BaseHeader(token, s.cnf.Sec.RefererNum))
	if err != nil {
		fmt.Println(fmt.Sprintf("token：%s 获取失败，错误信息：%s", token, err.Error()))
		return ""
	}

	rk := new(getRkResp)
	if err = json.Unmarshal(res, rk); err != nil {
		return ""
	}
	fmt.Println(fmt.Sprintf("token：%s 获取成功：RK：%s", token, rk.Data.Rk))
	return rk.Data.Rk
}

func (s *Seckiller) Result() {
	for _, ts := range s.tokens {
		s.result(ts.Token)
	}
}

func (s *Seckiller) result(token string) {
	// https://mengniu-apig.xiaoyisz.com/mengniu-world-cup/mp/api/user/goods/list?timestamp=1670377607905&nonce=eVB2lLrofyyUhB2v&signature=43010893B85E257D84BC0B1056E787BB&page=1&pageSize=100
	// https://mengniu-apig.xiaoyisz.com/mengniu-world-cup/mp/api/user/goods/list?timestamp=1671246424470&nonce=HJHgsW47yHrzrqVX&signature=8E97B1B8F8C0E6271D179DC9C9A9B976&page=1&pageSize=100
	timestamp := utilx.GetTimestamp()
	nonce := utilx.GetNonce()
	sign := utilx.GetSign(timestamp, nonce, s.cnf.Sec.ClientKey, s.cnf.Sec.ClientSecret)
	url := fmt.Sprintf("%s/%s?timestamp=%s&nonce=%s&signature=%s&page=1&pageSize=100", s.cnf.Sec.Domain, "mp/api/user/goods/list", timestamp, nonce, sign)

	res, err := utilx.HttpGet(url, utilx.BaseHeader(token, s.cnf.Sec.RefererNum))
	if err != nil {
		fmt.Println("获取抢购结果失败，err= ", err)
		return
	}
	s.log(fmt.Sprintf("token == %s, 查询结果  == %s", token, string(res)))

}

// WithLog 写入日志
func WithLog() Option {
	return func(s *Seckiller) {
		s.writeLog = true
		s.writePath = s.cnf.Custom.LogPath
	}
}

// WithJsonId 获取jsonId
func WithJsonId() Option {
	return func(s *Seckiller) {
		// 获取jsonId
		jsonId, rewardNum := s.getJsonId(0)
		if jsonId == "" {
			return
		}
		s.jsonId = jsonId
		s.rewardNum = rewardNum
	}
}

// WithRk 获取rk，筛选过期的token
func WithRk() Option {
	return func(s *Seckiller) {
		// 获取rk 过期的token不计入
		var ts []tokens
		for _, token := range s.cnf.Custom.Tokens {
			rk := s.getRk(token)
			if rk == "" {
				continue
			}
			ts = append(ts, tokens{
				Token: token,
				Rk:    rk,
			})
		}
		s.tokens = ts
	}
}

func seckillSign(rk, timestamp, desKey string) (string, string) {
	k := utilx.DesRK(rk, desKey)
	reqId := utilx.RandStr(32)
	sign := fmt.Sprintf("requestId=%s&timestamp=%s&key=%s", reqId, timestamp, k)
	return strings.ToUpper(utilx.MD5(sign)), reqId
}
