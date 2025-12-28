package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-820/transactions/internal/mode/api"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Run gRPC API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			return api.Run(ctx)
		},
	}

	rootCmd.AddCommand(cmd)
}
