package es

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

//	将日志数据写入Elasticsearch

type ESClient struct {
	client      *elastic.Client
	index       string
	logDataChan chan interface{}
}

var (
	esClient = &ESClient{}
)

func Init(addr, index string, goroutineNum, maxSize int) (err error) {
	client, err := elastic.NewClient(elastic.SetURL("http://" + addr))
	if err != nil {
		// Handle error
		panic(err)
	}
	esClient.client = client
	esClient.index = index
	esClient.logDataChan = make(chan interface{}, maxSize)
	fmt.Println("connect to es success")
	// 从通道中取数据
	for i := 0; i < goroutineNum; i++ {
		go sendToES()
	}
	return
}

func sendToES() {
	for m1 := range esClient.logDataChan {
		// msg, err := json.Marshal(m1)
		// if err != nil {
		// 	fmt.Printf("marshal m1 failed, err:%v\n", err)
		// 	continue
		// }
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(m1).
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// 通过一个首字母大写的函数从包外接受msg，发送到chan中
func PutLogData(msg interface{}) {
	esClient.logDataChan <- msg
}
