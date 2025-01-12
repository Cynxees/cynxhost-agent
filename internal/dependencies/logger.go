package dependencies

import (
	// "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	// elasticLog "gopkg.in/sohlich/elogrus.v7"
)

func NewLogger(config *Config) *logrus.Logger {
	esLogger := logrus.New()
	esLogger.SetFormatter(&ecslogrus.Formatter{})
	esLogger.SetLevel(logrus.InfoLevel)

	// client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	// if err != nil {
	// 	esLogger.Fatalf("Failed to create elasticsearch client: %v", err)
	// }

	// hook, err := elasticLog.NewAsyncElasticHook(client, config.App.Name, logrus.InfoLevel, "cynxhost-logs")
	// if err != nil {
	// 	logrus.Fatalf("Failed to create Elasticsearch hook: %v", err)
	// }
	// esLogger.AddHook(hook)

	return esLogger
}
