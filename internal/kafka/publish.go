package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

func PublishMessage(ctx context.Context, brokers []string, topic, key, value string) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return fmt.Errorf("publish: failed to create client: %w", err)
	}

	defer client.Close()

	var recordKey []byte
	if key != "" {
		recordKey = []byte(key)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   recordKey,
		Value: []byte(value),
	}

	if err := client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("publish: failed to produce record: %w", err)
	}

	fmt.Printf("published to %s (partition=%d offset=%d)\n",
		record.Topic, record.Partition, record.Offset)

	return nil
}
