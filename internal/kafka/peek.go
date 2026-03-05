package kafka

import (
	"context"
	"fmt"
)

func PeekMessages(client *Client, topic string, n int) error {

	ctx := context.Background()

	client.Kafka.AddConsumeTopics(topic)

	count := 0

	for {
		fetches := client.Kafka.PollFetches(ctx)

		iter := fetches.RecordIter()

		for !iter.Done() {
			record := iter.Next()

			fmt.Printf("offset=%d partition=%d\n",
				record.Offset,
				record.Partition,
			)

			fmt.Println(string(record.Value))
			fmt.Println()

			count++

			if count >= n {
				return nil
			}
		}
	}
}
