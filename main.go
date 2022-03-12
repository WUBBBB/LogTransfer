package main

import (
	"fmt"
	"main/es"
	"main/kafka"
	"main/model"

	"gopkg.in/ini.v1"
)

// 从kafka消费日志数据，写入ES

func main() {
	// 1.加载配置文件
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, "./config/logtransfer.ini")
	if err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		panic(err)
	}
	fmt.Println("load config success")
	// 2.连接kafka
	err = kafka.Init([]string{cfg.KafkaConf.Adderss}, cfg.KafkaConf.Topic)
	if err != nil {
		fmt.Printf("connected to kafka failed, err:%v\n", err)
		panic(err)
	}
	fmt.Println("connect to kafka success")
	// 3.连接ES
	err = es.Init(cfg.ESConf.Adderss, cfg.ESConf.Index, cfg.ESConf.GoNum, cfg.ESConf.MaxSize)
	if err != nil {
		fmt.Printf("Init es failed, err:%v\n", err)
		panic(err)
	}
	fmt.Println("Init ES success")
	// 在这儿停顿！让程序不退出，用select不用for是因为select不占CPU
	select {}
}
