package transactions

import "time"

type WalletUpdateEvent struct {
	EventID    string
	OccurredAt time.Time
	UserID     string
	Wallets    []WalletRef
}

type WalletRef struct {
	Chain   string
	Address string
	Label   string
}
