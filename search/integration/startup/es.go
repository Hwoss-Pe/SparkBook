package startup

import (
	"Webook/search/repository/dao"
	"github.com/olivere/elastic/v7"
	"log"
	"time"
)

func InitESClient() *elastic.Client {
	const timeout = 10 * time.Second
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeoutStartup(timeout),
		elastic.SetTraceLog(log.Default()),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	err = dao.InitES(client)
	if err != nil {
		panic(err)
	}
	return client
}
