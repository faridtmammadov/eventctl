package cmd

import (
	"github.com/faridtmammadov/eventctl/internal/kafka"
	"github.com/spf13/cobra"
)

var num int

var peekCmd = &cobra.Command{
	Use:   "peek [topic]",
	Short: "Peek messages from a topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		client := kafka.NewClient("localhost:9092")

		return kafka.PeekMessages(client, topic, num)
	},
}

func init() {
	peekCmd.Flags().IntVarP(&num, "number", "n", 5, "number of messages")
	rootCmd.AddCommand(peekCmd)
}
