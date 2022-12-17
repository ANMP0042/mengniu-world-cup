/**
 * @Author: YMBoom
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2022/12/05 17:11
 */
package main

import (
	"fmt"
	"mengniu/seckill"
	"net/http"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	sler, err := seckill.NewSeckiller(seckill.WithLog(), seckill.WithJsonId(), seckill.WithRk())
	if err != nil {
		return
	}

	go sler.Seckill()

	fmt.Println("listenAndServe on", ":8000")
	//sler.Result()
	http.ListenAndServe(":8000", nil)
}
