package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/project-820/transactions/internal/adapters/out/onchain/httpx"
	"github.com/project-820/transactions/internal/core/usecase"
)

var _ usecase.OnchainClient = (*Client)(nil)

const transferTopic0 = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" // keccak("Transfer(address,address,uint256)")

type ClientParams struct {
	RPCUrl          string
	Doer            httpx.Doer
	MaxBlocksPerRun uint64
}

type Client struct {
	rpcURL          string
	doer            httpx.Doer
	maxBlocksPerRun uint64
}

func NewEVMClient(params ClientParams) *Client {
	if params.MaxBlocksPerRun == 0 {
		params.MaxBlocksPerRun = 2000
	}

	return &Client{
		rpcURL:          params.RPCUrl,
		doer:            params.Doer,
		maxBlocksPerRun: params.MaxBlocksPerRun,
	}
}

func (c *Client) ListTransactions(
	ctx context.Context,
	req usecase.OnchainTxListRequest,
) (usecase.OnchainTxListResult, error) {
	addr := normalizeEVMAddress(req.Address)
	if addr == "" {
		return usecase.OnchainTxListResult{}, fmt.Errorf("invalid evm address: %q", req.Address)
	}

	last := uint64(0)
	if strings.TrimSpace(req.Cursor) != "" {
		var ok bool
		last, ok = parseUint(req.Cursor)
		if !ok {
			return usecase.OnchainTxListResult{}, fmt.Errorf("invalid cursor (expected decimal block): %q", req.Cursor)
		}
	}

	latestHex, err := c.blockNumber(ctx)
	if err != nil {
		return usecase.OnchainTxListResult{}, err
	}
	latest := mustHexToUint(latestHex)

	from := last + 1
	to := latest
	if from > to {
		return usecase.OnchainTxListResult{NextCursor: fmt.Sprintf("%d", last), HasMore: false}, nil
	}

	if to-from+1 > c.maxBlocksPerRun {
		to = from + c.maxBlocksPerRun - 1
	}

	logsIn, err := c.getTransferLogs(ctx, from, to, "", addr)
	if err != nil {
		return usecase.OnchainTxListResult{}, err
	}
	logsOut, err := c.getTransferLogs(ctx, from, to, addr, "")
	if err != nil {
		return usecase.OnchainTxListResult{}, err
	}

	// merge + dedupe by (txhash, logindex)
	seen := make(map[string]struct{}, len(logsIn)+len(logsOut))
	out := make([]usecase.OnchainTx, 0, len(logsIn)+len(logsOut))

	// Время: без дополнительного запроса к блоку будет неизвестно.
	// Для MVP можно оставить OccurredAt=0 и восстановить позже.
	// Если тебе обязательно — добавишь кеш blockTimestampByNumber + eth_getBlockByNumber.
	for _, lg := range append(logsIn, logsOut...) {
		k := strings.ToLower(lg.TransactionHash) + ":" + strings.ToLower(lg.LogIndex)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}

		val, err := hexUint256ToDecimal(lg.Data)
		if err != nil {
			return usecase.OnchainTxListResult{}, fmt.Errorf("bad log value: %w", err)
		}

		blockNum := mustHexToUint(lg.BlockNumber)
		logIndex := int32(mustHexToUint(lg.LogIndex))

		out = append(out, usecase.OnchainTx{
			TxHash:     strings.ToLower(lg.TransactionHash),
			LogIndex:   logIndex,
			BlockNum:   int64(blockNum),
			OccurredAt: time.Time{},                                                             // TODO optional: fetch timestamps
			AssetRef:   "evm:" + strings.ToLower(req.Chain) + ":" + strings.ToLower(lg.Address), // token contract
			Qty:        val,
			Meta:       nil,
		})
	}

	nextCursor := fmt.Sprintf("%d", to)
	hasMore := to < latest

	return usecase.OnchainTxListResult{
		Transactions: out,
		NextCursor:   nextCursor,
		HasMore:      hasMore,
	}, nil
}

func (c *Client) blockNumber(ctx context.Context) (string, error) {
	var out string
	if err := httpx.CallRPC(ctx, c.doer, c.rpcURL, 1, "eth_blockNumber", []any{}, &out); err != nil {
		return "", err
	}

	return out, nil
}

func (c *Client) getTransferLogs(ctx context.Context, fromBlock, toBlock uint64, fromAddr, toAddr string) ([]rpcLog, error) {
	filter := map[string]any{
		"fromBlock": fmt.Sprintf("0x%x", fromBlock),
		"toBlock":   fmt.Sprintf("0x%x", toBlock),
		"topics":    buildTopics(fromAddr, toAddr),
	}
	var logs []rpcLog
	if err := httpx.CallRPC(ctx, c.doer, c.rpcURL, 2, "eth_getLogs", []any{filter}, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// topics: [transferSig, from, to]
func buildTopics(fromAddr, toAddr string) []any {
	topics := make([]any, 3)
	topics[0] = transferTopic0
	if fromAddr != "" {
		topics[1] = evmTopicAddress(fromAddr)
	}
	if toAddr != "" {
		topics[2] = evmTopicAddress(toAddr)
	}
	return topics
}

// address must be 0x + 40 hex
func evmTopicAddress(addr string) string {
	a := strings.TrimPrefix(strings.ToLower(addr), "0x")
	return "0x" + strings.Repeat("0", 24*2) + a
}

func normalizeEVMAddress(a string) string {
	s := strings.ToLower(strings.TrimSpace(a))
	if strings.HasPrefix(s, "0x") && len(s) == 42 {
		_, err := hex.DecodeString(s[2:])
		if err == nil {
			return s
		}
	}
	return ""
}

func mustHexToUint(hexStr string) uint64 {
	x, _ := new(big.Int).SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	return x.Uint64()
}

func hexUint256ToDecimal(hexStr string) (string, error) {
	s := strings.TrimPrefix(hexStr, "0x")
	if s == "" {
		return "0", nil
	}
	n := new(big.Int)
	_, ok := n.SetString(s, 16)
	if !ok {
		return "", fmt.Errorf("bad hex: %q", hexStr)
	}
	return n.String(), nil
}

func parseUint(s string) (uint64, bool) {
	n := new(big.Int)
	_, ok := n.SetString(strings.TrimSpace(s), 10)
	if !ok {
		return 0, false
	}
	return n.Uint64(), true
}
