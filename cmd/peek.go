package cmd

import (
	"fmt"

	"github.com/faridtmammadov/eventctl/internal/config"
	"github.com/faridtmammadov/eventctl/internal/kafka"
	"github.com/spf13/cobra"
)

var num int

var peekCmd = &cobra.Command{
	Use:   "peek [topic]",
	Short: "Peek messages from a topic",
	Long:  `Peek prints the most recent messages from a topic.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		topic := args[0]

		if num <= 0 {
			return fmt.Errorf("number must be greater than 0")
		}

		conn, err := config.BuildConnectionConfig(cfg, connName, brokerOverride)
		if err != nil {
			return fmt.Errorf("build connection config: %w", err)
		}

		client := kafka.NewClient(conn.Brokers...)
		defer client.Close()

		return kafka.PeekMessages(ctx, client, topic, num)
	},
}

func init() {
	peekCmd.Flags().IntVarP(&num, "number", "n", 1, "number of messages to print")
	rootCmd.AddCommand(peekCmd)
}
