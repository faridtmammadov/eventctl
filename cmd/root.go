package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/faridtmammadov/eventctl/internal/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var (
	connName       string
	brokerOverride string
)

var rootCmd = &cobra.Command{
	Use:          "eventctl",
	Short:        "Event debugging CLI",
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&connName,
		"connection",
		"c",
		"",
		"named connection from ~/.eventctl/config.yaml",
	)

	rootCmd.PersistentFlags().StringVarP(
		&brokerOverride,
		"broker",
		"b",
		"",
		"broker address(es), comma-separated (e.g. host:9092,host2:9092)",
	)
}

func ExecuteContext(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
