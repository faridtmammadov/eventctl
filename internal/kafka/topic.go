package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kadm"
)

func CreateTopic(ctx context.Context, client *Client, name string, partitions int32, replicationFactor int16) error {
	adm := kadm.NewClient(client.Kafka)

	resp, err := adm.CreateTopics(ctx, partitions, replicationFactor, nil, name)
	if err != nil {
		return fmt.Errorf("topic create: request failed: %w", err)
	}

	for _, detail := range resp {
		if detail.Err != nil {
			return fmt.Errorf("topic create: broker rejected %q: %w", detail.Topic, detail.Err)
		}
		fmt.Printf("created topic %q (partitions=%d replication-factor=%d)\n",
			detail.Topic, partitions, replicationFactor)
	}

	return nil
}
