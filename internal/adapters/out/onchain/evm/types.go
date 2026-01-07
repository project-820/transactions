package evm

type rpcLog struct {
	Address         string   `json:"address"`
	Topics          []string `json:"topics"`
	Data            string   `json:"data"`
	BlockNumber     string   `json:"blockNumber"`
	TransactionHash string   `json:"transactionHash"`
	LogIndex        string   `json:"logIndex"`
}

// {
//     "jsonrpc": "2.0",
//     "id": "14",
//     "result": [
//         {
//             "address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//             "topics": [
//                 "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c",
//                 "0x000000000000000000000000111111125421ca6dc452d289314280a0f8842a65"
//             ],
//             "data": "0x00000000000000000000000000000000000000000000000001ede17395e65f21",
//             "blockNumber": "0x170fc37",
//             "transactionHash": "0x09cc18ff4aa49421467b53b97064efc7b4845ff846908c22d8bbff9a1600e5c8",
//             "transactionIndex": "0x0",
//             "blockHash": "0x4487a7cfea9f37d02c45d323adcbb81946d84cedc01bd6fcaf5e6869e1dfa2bc",
//             "blockTimestamp": "0x0",
//             "logIndex": "0x0",
//             "removed": false
//         }
// 	]
// }
