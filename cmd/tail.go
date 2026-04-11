package cmd

import (
	"fmt"

	"github.com/faridtmammadov/eventctl/internal/config"
	"github.com/faridtmammadov/eventctl/internal/kafka"
	"github.com/spf13/cobra"
)

var num int

var tailCmd = &cobra.Command{
	Use:   "tail [topic]",
	Short: "Tail messages from a topic",
	Long:  `Tail prints the most recent messages from a topic.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		topic := args[0]

		if num <= 0 {
			return fmt.Errorf("tail: --number must be greater than 0")
		}

		conn, err := config.BuildConnectionConfig(cfg, connName, brokerOverride)
		if err != nil {
			return fmt.Errorf("build connection config: %w", err)
		}

		return kafka.Tail(ctx, conn.Brokers, topic, num)
	},
}

func init() {
	tailCmd.Flags().IntVarP(&num, "number", "n", 1, "number of messages to print")
	rootCmd.AddCommand(tailCmd)
}
