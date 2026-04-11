package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

func PeekMessages(ctx context.Context, brokers []string, topic string, n int) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	adm := kadm.NewClient(client)
	listed, err := adm.ListEndOffsets(ctx, topic)
	if err != nil {
		return fmt.Errorf("peek: failed to list end offsets: %w", err)
	}

	pending := make(map[int32]int64)
	listed.Each(func(lo kadm.ListedOffset) {
		if lo.Offset > 0 {
			pending[lo.Partition] = lo.Offset
		}
	})

	if len(pending) == 0 {
		fmt.Println("(no messages)")
		return nil
	}

	client.AddConsumeTopics(topic)

	count := 0

	for count < n && len(pending) > 0 {
		fetches := client.PollFetches(ctx)

		if err := fetches.Err(); err != nil {
			return fmt.Errorf("peek: fetch error: %w", err)
		}

		if fetches.NumRecords() == 0 {
			break
		}

		iter := fetches.RecordIter()
		for !iter.Done() && count < n {
			record := iter.Next()

			fmt.Printf("offset=%d partition=%d\n", record.Offset, record.Partition)
			fmt.Println(string(record.Value))
			fmt.Println()

			count++

			if end, ok := pending[record.Partition]; ok && record.Offset+1 >= end {
				delete(pending, record.Partition)
			}
		}
	}

	return nil
}
