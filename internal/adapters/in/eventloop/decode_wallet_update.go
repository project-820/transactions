package eventloop

import (
	"encoding/json"
	"time"

	"github.com/project-820/transactions/internal/core/usecase"
)

type walletUpdate struct {
	EventID    string    `json:"event_id"`
	OccurredAt time.Time `json:"occurred_at"`
	UserID     string    `json:"user_id"`
	Wallets    []struct {
		Chain   string `json:"chain"`
		Address string `json:"address"`
		Label   string `json:"label,omitempty"`
	} `json:"wallets"`
}

func decodeWalletUpdate(payload []byte) (usecase.WalletUpdateEvent, error) {
	var walletUpdate walletUpdate
	if err := json.Unmarshal(payload, &walletUpdate); err != nil {
		return usecase.WalletUpdateEvent{}, err
	}

	walletUpdateEvent := usecase.WalletUpdateEvent{
		EventID:    walletUpdate.EventID,
		OccurredAt: walletUpdate.OccurredAt,
		UserID:     walletUpdate.UserID,
		Wallets:    make([]usecase.WalletRef, 0, len(walletUpdate.Wallets)),
	}

	for _, wallet := range walletUpdate.Wallets {
		walletUpdateEvent.Wallets = append(walletUpdateEvent.Wallets, usecase.WalletRef{
			Chain:   wallet.Chain,
			Address: wallet.Address,
			Label:   wallet.Label,
		})
	}

	return walletUpdateEvent, nil
}
