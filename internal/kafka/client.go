package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	Kafka *kgo.Client
}

func NewClient(brokers string) *Client {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers),
	)
	if err != nil {
		panic(err)
	}

	return &Client{Kafka: cl}
}
