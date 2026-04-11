package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

func PublishMessage(ctx context.Context, client *Client, topic, key, value string) error {
	var recordKey []byte
	if key != "" {
		recordKey = []byte(key)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   recordKey,
		Value: []byte(value),
	}

	if err := client.Kafka.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("publish: failed to produce record: %w", err)
	}

	fmt.Printf("published to %s (partition=%d offset=%d)\n",
		record.Topic, record.Partition, record.Offset)

	return nil
}
