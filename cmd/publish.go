package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/faridtmammadov/eventctl/internal/config"
	"github.com/faridtmammadov/eventctl/internal/kafka"
	"github.com/spf13/cobra"
)

var (
	publishMessage string
	publishKey     string
)

var publishCmd = &cobra.Command{
	Use:   "publish [topic]",
	Short: "Publish a message to a topic",
	Long: `Publish sends a single message to a topic.

The message value is read from --message or from stdin when --message is omitted.

Examples:
  # Inline value
  eventctl publish orders --message '{"id": 1, "status": "new"}'

  # With an explicit key
  eventctl publish orders --key order-123 --message '{"id": 1}'

  # Pipe from stdin
  echo '{"id": 1}' | eventctl publish orders`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		topic := args[0]

		value, err := resolveMessage(cmd, publishMessage)
		if err != nil {
			return err
		}

		conn, err := config.BuildConnectionConfig(cfg, connName, brokerOverride)
		if err != nil {
			return err
		}

		return kafka.PublishMessage(ctx, conn.Brokers, topic, publishKey, value)
	},
}

func resolveMessage(cmd *cobra.Command, flagValue string) (string, error) {
	if cmd.Flags().Changed("message") {
		if flagValue == "" {
			return "", errors.New("publish: --message cannot be empty")
		}
		return flagValue, nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("publish: failed to stat stdin: %w", err)
	}

	if (fi.Mode() & os.ModeCharDevice) != 0 {
		return "", errors.New("publish: message is empty; provide --message or pipe input via stdin")
	}

	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("publish: failed to read stdin: %w", err)
	}

	value := strings.TrimRight(string(bytes), "\r\n")
	if value == "" {
		return "", errors.New("publish: message is empty after reading stdin")
	}

	return value, nil
}

func init() {
	publishCmd.Flags().StringVarP(&publishMessage, "message", "m", "", "message value to publish (reads stdin when omitted)")
	publishCmd.Flags().StringVarP(&publishKey, "key", "k", "", "record key (optional; Kafka assigns partition automatically when empty)")
	rootCmd.AddCommand(publishCmd)
}
