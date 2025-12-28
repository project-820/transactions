package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-820/transactions/internal/mode/worker"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Run worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			return worker.Run(ctx)
		},
	}

	rootCmd.AddCommand(cmd)
}
