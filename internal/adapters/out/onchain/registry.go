package onchain

import (
	"strings"

	"github.com/project-820/transactions/internal/core/usecase"
)

var _ usecase.OnchainResolver = (*Registry)(nil)

type Registry struct {
	evm usecase.OnchainClient
}

func NewRegistry(evm usecase.OnchainClient) *Registry {
	return &Registry{
		evm: evm,
	}
}

func (r *Registry) ForChain(chain string) (usecase.OnchainClient, bool) {
	switch strings.ToLower(strings.TrimSpace(chain)) {
	case "eth", "ethereum", "bsc", "bnb", "polygon", "matic", "arbitrum", "arb", "optimism", "op", "base", "avax", "avalanche":
		return r.evm, r.evm != nil
	default:
		return nil, false
	}
}
