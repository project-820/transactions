package ports

import (
	"context"

	"github.com/project-820/transactions/internal/services/entity"
)

type OutboxRepository interface {
	AddMessages(ctx context.Context, msgs []entity.OutboxMessage) error
}
