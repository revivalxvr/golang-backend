package env

import "os"

type Consumer struct {
	KafkaBrokers       []string
	KafkaTopic         string
	KafkaGroupID       string
	KafkaDLQTopic      string
	ElasticsearchUrl   []string
	ElasticsearchIndex string
}
type Config struct {
	RestfulLogConsumer  *Consumer
	RabbitMqLogConsumer *Consumer
	GrpcLogConsumer     *Consumer
}

var Conf *Config

func Load() {
	kafkaAddr1 := os.Getenv("KAFKA_ADDR_1")
	kafkaAddr2 := os.Getenv("KAFKA_ADDR_2")
	kafkaAddr3 := os.Getenv("KAFKA_ADDR_3")

	elasticAddr1 := os.Getenv("ELASTICSEARCH_ADDR_1")
	elasticAddr2 := os.Getenv("ELASTICSEARCH_ADDR_2")
	elasticAddr3 := os.Getenv("ELASTICSEARCH_ADDR_3")

	restfulLogConsumerConf := new(Consumer)
	restfulLogConsumerConf.KafkaBrokers = []string{kafkaAddr1, kafkaAddr2, kafkaAddr3}
	restfulLogConsumerConf.KafkaTopic = "restful-log"
	restfulLogConsumerConf.KafkaGroupID = "restful-log-group"
	restfulLogConsumerConf.KafkaDLQTopic = "restful-log-dlq"
	restfulLogConsumerConf.ElasticsearchUrl = []string{elasticAddr1, elasticAddr2, elasticAddr3}
	restfulLogConsumerConf.ElasticsearchIndex = "restful-log"

	rabbitMQLogConsumerConf := new(Consumer)
	rabbitMQLogConsumerConf.KafkaBrokers = []string{kafkaAddr1, kafkaAddr2, kafkaAddr3}
	rabbitMQLogConsumerConf.KafkaTopic = "rabbitmq-log"
	rabbitMQLogConsumerConf.KafkaGroupID = "rabbitmq-log-group"
	rabbitMQLogConsumerConf.KafkaDLQTopic = "rabbitmq-log-dlq"
	rabbitMQLogConsumerConf.ElasticsearchUrl = []string{elasticAddr1, elasticAddr2, elasticAddr3}
	rabbitMQLogConsumerConf.ElasticsearchIndex = "rabbitmq-log"

	grpcLogConsumerConf := new(Consumer)
	grpcLogConsumerConf.KafkaBrokers = []string{kafkaAddr1, kafkaAddr2, kafkaAddr3}
	grpcLogConsumerConf.KafkaTopic = "grpc-log"
	grpcLogConsumerConf.KafkaGroupID = "grpc-log-group"
	grpcLogConsumerConf.KafkaDLQTopic = "grpc-log-dlq"
	grpcLogConsumerConf.ElasticsearchUrl = []string{elasticAddr1, elasticAddr2, elasticAddr3}
	grpcLogConsumerConf.ElasticsearchIndex = "grpc-log"

	Conf = &Config{
		RestfulLogConsumer:  restfulLogConsumerConf,
		RabbitMqLogConsumer: rabbitMQLogConsumerConf,
		GrpcLogConsumer:     grpcLogConsumerConf,
	}
}
