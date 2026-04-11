package cmd

import (
	"github.com/faridtmammadov/eventctl/internal/config"
	"github.com/faridtmammadov/eventctl/internal/kafka"
	"github.com/spf13/cobra"
)

var (
	topicPartitions        int32
	topicReplicationFactor int16
)

var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Manage Kafka topics",
}

var topicCreateCmd = &cobra.Command{
	Use:   "create [topic]",
	Short: "Create a topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		conn, err := config.BuildConnectionConfig(cfg, connName, brokerOverride)
		if err != nil {
			return err
		}

		client := kafka.NewClient(conn.Brokers...)
		defer client.Close()

		return kafka.CreateTopic(ctx, client, name, topicPartitions, topicReplicationFactor)
	},
}

func init() {
	topicCreateCmd.Flags().Int32VarP(&topicPartitions, "partitions", "p", 1, "number of partitions")
	topicCreateCmd.Flags().Int16VarP(&topicReplicationFactor, "replication-factor", "r", 1, "replication factor")

	topicCmd.AddCommand(topicCreateCmd)
	rootCmd.AddCommand(topicCmd)
}
