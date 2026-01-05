package evm

type rpcLog struct {
	Address         string   `json:"address"`     // token contract
	Topics          []string `json:"topics"`      // [TransferSig, from, to]
	Data            string   `json:"data"`        // value (uint256) hex
	BlockNumber     string   `json:"blockNumber"` // hex
	TransactionHash string   `json:"transactionHash"`
	LogIndex        string   `json:"logIndex"` // hex
}
