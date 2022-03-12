package kafka

import (
	"encoding/json"
	"fmt"
	"main/es"

	"github.com/Shopify/sarama"
)

// 初始化kafka连接
// 从kafka里面取出日志数据

func Init(addr []string, topic string) (err error) {
	// 创建新的消费者
	consumer, err := sarama.NewConsumer(addr, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	// 拿到指定topic下面的所有分区列表
	partitionList, err := consumer.Partitions(topic) //根据参数topic收到所有分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	// var wg sync.WaitGroup   	// wg相关内容需要注释掉，否则程序无法退出导致main的后续内容无法执行
	for partition := range partitionList { //遍历所有的分区
		//针对每个分区创建一个对应的分区消费者
		pc, err2 := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err2 != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n",
				partition, err2)
			return
		}
		// defer pc.AsyncClose()
		//异步从每个分区消费信息
		// wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() { //有个问题，这个Mseeages是个channel吧，这玩意只进不出的吗？！迟早会装满的吧
				// logDataChan <- msg //	为了将同步流程异步化，所以将取出的日志数据先放到通道中
				//	改成了下方的代码
				var m1 map[string]interface{}
				err = json.Unmarshal(msg.Value, &m1)
				if err != nil {
					fmt.Printf("unmarshal msg failed, err:%v\n", err)
					continue
				}
				es.PutLogData(m1)
			}
		}(pc)
	}
	// wg.Wait()
	return
}
